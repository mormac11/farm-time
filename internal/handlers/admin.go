package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"farm-time/internal/auth"
	"farm-time/internal/models"
)

// Admin handlers - require admin privileges

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	user := auth.GetUserFromContext(r.Context())
	if user == nil || !user.IsAdmin {
		h.respondError(w, http.StatusForbidden, "Admin access required")
		return
	}

	users, err := h.db.ListAllUsers(r.Context())
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to list users")
		return
	}
	if users == nil {
		users = []models.User{}
	}
	h.respondJSON(w, http.StatusOK, users)
}

func (h *Handler) UpdateUserPermissions(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	currentUser := auth.GetUserFromContext(r.Context())
	if currentUser == nil || !currentUser.IsAdmin {
		h.respondError(w, http.StatusForbidden, "Admin access required")
		return
	}

	userID := chi.URLParam(r, "userId")

	var req models.UpdateUserPermissionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.db.UpdateUserPermissions(r.Context(), userID, req)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "User not found")
		return
	}

	h.respondJSON(w, http.StatusOK, user)
}
