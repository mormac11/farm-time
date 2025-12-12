package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"farm-time/internal/db"
	"farm-time/internal/models"
)

type contextKey string

const UserContextKey contextKey = "user"

type Handler struct {
	db           *db.DB
	oauth2Config *oauth2.Config
	sessionTTL   time.Duration
	cookieName   string
	secureCookie bool
}

func NewHandler(database *db.DB) *Handler {
	return &Handler{
		db: database,
		oauth2Config: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
		sessionTTL:   7 * 24 * time.Hour, // 7 days
		cookieName:   "farm_session",
		secureCookie: os.Getenv("ENV") == "production",
	}
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (h *Handler) setSessionCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.cookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.secureCookie,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.sessionTTL.Seconds()),
	})
}

func (h *Handler) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   h.secureCookie,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

// GoogleLogin redirects the user to Google's OAuth consent page
func (h *Handler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state, err := generateSessionID()
	if err != nil {
		http.Error(w, `{"error":"Failed to generate state"}`, http.StatusInternalServerError)
		return
	}

	// Store state in a short-lived cookie for CSRF protection
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.secureCookie,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300, // 5 minutes
	})

	url := h.oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallback handles the OAuth callback from Google
func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Verify state for CSRF protection
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
		http.Error(w, `{"error":"Invalid state"}`, http.StatusBadRequest)
		return
	}

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Check for error from Google
	if errMsg := r.URL.Query().Get("error"); errMsg != "" {
		http.Redirect(w, r, "/?error="+errMsg, http.StatusTemporaryRedirect)
		return
	}

	// Exchange code for token
	code := r.URL.Query().Get("code")
	token, err := h.oauth2Config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, `{"error":"Failed to exchange token"}`, http.StatusInternalServerError)
		return
	}

	// Get user info from Google
	client := h.oauth2Config.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, `{"error":"Failed to get user info"}`, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		http.Error(w, `{"error":"Failed to decode user info"}`, http.StatusInternalServerError)
		return
	}

	// Create or update user in database
	user, err := h.db.CreateOrUpdateUser(r.Context(), googleUser.ID, googleUser.Email, googleUser.Name, googleUser.Picture)
	if err != nil {
		http.Error(w, `{"error":"Failed to save user"}`, http.StatusInternalServerError)
		return
	}

	// Create session
	sessionID, err := generateSessionID()
	if err != nil {
		http.Error(w, `{"error":"Failed to create session"}`, http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(h.sessionTTL)
	if err := h.db.CreateSession(r.Context(), sessionID, user.ID, expiresAt); err != nil {
		http.Error(w, `{"error":"Failed to save session"}`, http.StatusInternalServerError)
		return
	}

	// Set session cookie
	h.setSessionCookie(w, sessionID)

	// Redirect to app
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// Logout clears the user's session
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(h.cookieName)
	if err == nil {
		h.db.DeleteSession(r.Context(), cookie.Value)
	}

	h.clearSessionCookie(w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "logged out"})
}

// Me returns the current user's info
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user, _ := h.GetUserFromRequest(r)

	w.Header().Set("Content-Type", "application/json")
	if user == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"user": nil})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"user": user})
}

// GetUserFromRequest extracts the user from the session cookie
func (h *Handler) GetUserFromRequest(r *http.Request) (*models.User, error) {
	cookie, err := r.Cookie(h.cookieName)
	if err != nil {
		return nil, err
	}

	session, err := h.db.GetSession(r.Context(), cookie.Value)
	if err != nil {
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		h.db.DeleteSession(r.Context(), session.ID)
		return nil, errors.New("session expired")
	}

	return h.db.GetUserByID(r.Context(), session.UserID)
}

// RequireAuth middleware returns 401 if the user is not authenticated
func (h *Handler) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := h.GetUserFromRequest(r)
		if err != nil || user == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext retrieves the user from the request context
func GetUserFromContext(ctx context.Context) *models.User {
	user, ok := ctx.Value(UserContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}
