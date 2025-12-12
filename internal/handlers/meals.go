package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"farm-time/internal/auth"
	"farm-time/internal/models"
)

// Meal handlers

func (h *Handler) ListMeals(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")

	meals, err := h.db.GetMealsWithItemsByEvent(r.Context(), eventID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to list meals")
		return
	}
	if meals == nil {
		meals = []models.MealWithItems{}
	}
	h.respondJSON(w, http.StatusOK, meals)
}

func (h *Handler) CreateMeal(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")

	var req models.CreateMealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		h.respondError(w, http.StatusBadRequest, "Name is required")
		return
	}
	if req.MealType == "" {
		h.respondError(w, http.StatusBadRequest, "Meal type is required")
		return
	}

	meal, err := h.db.CreateMeal(r.Context(), eventID, req)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create meal")
		return
	}

	h.respondJSON(w, http.StatusCreated, meal)
}

func (h *Handler) UpdateMeal(w http.ResponseWriter, r *http.Request) {
	mealID := chi.URLParam(r, "mealId")

	var req models.UpdateMealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	meal, err := h.db.UpdateMeal(r.Context(), mealID, req)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "Meal not found")
		return
	}

	h.respondJSON(w, http.StatusOK, meal)
}

func (h *Handler) DeleteMeal(w http.ResponseWriter, r *http.Request) {
	mealID := chi.URLParam(r, "mealId")

	if err := h.db.DeleteMeal(r.Context(), mealID); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to delete meal")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// MealItem handlers

func (h *Handler) AddMealItem(w http.ResponseWriter, r *http.Request) {
	mealID := chi.URLParam(r, "mealId")

	var req models.CreateMealItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		h.respondError(w, http.StatusBadRequest, "Name is required")
		return
	}

	item, err := h.db.CreateMealItem(r.Context(), mealID, req)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to add meal item")
		return
	}

	h.respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateMealItem(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "itemId")

	var req models.UpdateMealItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.db.UpdateMealItem(r.Context(), itemID, req)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "Meal item not found")
		return
	}

	h.respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteMealItem(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "itemId")

	if err := h.db.DeleteMealItem(r.Context(), itemID); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to delete meal item")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// MealSignup handlers

func (h *Handler) SignupForItem(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	itemID := chi.URLParam(r, "itemId")

	// Get current user from context
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.CreateMealSignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Allow empty body for signups without notes
		req = models.CreateMealSignupRequest{}
	}

	// Create the signup
	signup, err := h.db.CreateMealSignup(r.Context(), itemID, user.ID, req)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to sign up for item")
		return
	}

	// Auto-add user as attendee if not already
	attendees, _ := h.db.GetAttendeesByEvent(r.Context(), eventID)
	isAttendee := false
	for _, a := range attendees {
		if a.Email == user.Email {
			isAttendee = true
			break
		}
	}
	if !isAttendee {
		h.db.CreateAttendee(r.Context(), eventID, models.CreateAttendeeRequest{
			Name:   user.Name,
			Email:  user.Email,
			Status: "attending",
		})
	}

	h.respondJSON(w, http.StatusCreated, signup)
}

func (h *Handler) RemoveSignup(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "itemId")

	// Get current user from context
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.db.DeleteMealSignup(r.Context(), itemID, user.ID); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to remove signup")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetEventWithMeals returns event with attendees and meals
func (h *Handler) GetEventWithMeals(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	event, err := h.db.GetEventWithMeals(r.Context(), id)
	if err != nil {
		log.Printf("GetEventWithMeals error for id %s: %v", id, err)
		h.respondError(w, http.StatusNotFound, "Event not found")
		return
	}

	h.respondJSON(w, http.StatusOK, event)
}
