package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"farm-time/internal/db"
	"farm-time/internal/models"
)

type Handler struct {
	db *db.DB
}

func New(database *db.DB) *Handler {
	return &Handler{db: database}
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}

// Event handlers

func (h *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.db.ListEvents(r.Context())
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to list events")
		return
	}
	if events == nil {
		events = []models.Event{}
	}
	h.respondJSON(w, http.StatusOK, events)
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" {
		h.respondError(w, http.StatusBadRequest, "Title is required")
		return
	}
	if req.StartTime.IsZero() || req.EndTime.IsZero() {
		h.respondError(w, http.StatusBadRequest, "Start time and end time are required")
		return
	}

	event, err := h.db.CreateEvent(r.Context(), req)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create event")
		return
	}

	// Auto-create meals based on event times
	h.autoCreateMeals(r.Context(), event.ID, req.StartTime, req.EndTime)

	h.respondJSON(w, http.StatusCreated, event)
}

func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	event, err := h.db.GetEventWithAttendees(r.Context(), id)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "Event not found")
		return
	}

	h.respondJSON(w, http.StatusOK, event)
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	event, err := h.db.UpdateEvent(r.Context(), id, req)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "Event not found")
		return
	}

	h.respondJSON(w, http.StatusOK, event)
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.db.DeleteEvent(r.Context(), id); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to delete event")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Attendee handlers

func (h *Handler) ListAttendees(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")

	attendees, err := h.db.GetAttendeesByEvent(r.Context(), eventID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to list attendees")
		return
	}
	if attendees == nil {
		attendees = []models.Attendee{}
	}
	h.respondJSON(w, http.StatusOK, attendees)
}

func (h *Handler) AddAttendee(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")

	var req models.CreateAttendeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		h.respondError(w, http.StatusBadRequest, "Name is required")
		return
	}
	if req.Email == "" {
		h.respondError(w, http.StatusBadRequest, "Email is required")
		return
	}

	attendee, err := h.db.CreateAttendee(r.Context(), eventID, req)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to add attendee")
		return
	}

	h.respondJSON(w, http.StatusCreated, attendee)
}

func (h *Handler) UpdateAttendee(w http.ResponseWriter, r *http.Request) {
	attendeeID := chi.URLParam(r, "attendeeId")

	var req models.UpdateAttendeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	attendee, err := h.db.UpdateAttendee(r.Context(), attendeeID, req)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "Attendee not found")
		return
	}

	h.respondJSON(w, http.StatusOK, attendee)
}

func (h *Handler) RemoveAttendee(w http.ResponseWriter, r *http.Request) {
	attendeeID := chi.URLParam(r, "attendeeId")

	if err := h.db.DeleteAttendee(r.Context(), attendeeID); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to remove attendee")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// autoCreateMeals creates lunch and dinner meals based on event times
// - Lunch: created if start time is at or before 11:00 AM
// - Dinner: created if end time is at or after 8:00 PM
// For multi-day events, creates meals for each day
func (h *Handler) autoCreateMeals(ctx context.Context, eventID string, startTime, endTime time.Time) {
	// Normalize to local dates
	startDate := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
	endDate := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, endTime.Location())

	// Iterate through each day of the event
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		dayName := d.Format("Monday")

		// Check if we should create lunch for this day
		// First day: check if start time is at or before 11 AM
		// Last day: check if end time is at or after 12 PM (noon)
		// Middle days: always create lunch
		createLunch := false
		if d.Equal(startDate) && d.Equal(endDate) {
			// Single day event - check start time
			if startTime.Hour() <= 11 {
				createLunch = true
			}
		} else if d.Equal(startDate) {
			// First day of multi-day - check start time
			if startTime.Hour() <= 11 {
				createLunch = true
			}
		} else if d.Equal(endDate) {
			// Last day - only create lunch if event goes past noon
			if endTime.Hour() >= 12 {
				createLunch = true
			}
		} else {
			// Middle day - always create lunch
			createLunch = true
		}

		if createLunch {
			mealName := "Lunch"
			if !startDate.Equal(endDate) {
				mealName = dayName + " Lunch"
			}
			h.db.CreateMeal(ctx, eventID, models.CreateMealRequest{
				Name:     mealName,
				MealType: "lunch",
				MealDate: &dateStr,
			})
		}

		// Check if we should create dinner for this day
		// Last day: check if end time is at or after 8 PM (20:00)
		// Other days: always create dinner
		createDinner := false
		if d.Equal(endDate) {
			// Last day - check end time
			if endTime.Hour() >= 20 {
				createDinner = true
			}
		} else {
			// Not the last day, create dinner
			createDinner = true
		}

		if createDinner {
			mealName := "Dinner"
			if !startDate.Equal(endDate) {
				mealName = dayName + " Dinner"
			}
			h.db.CreateMeal(ctx, eventID, models.CreateMealRequest{
				Name:     mealName,
				MealType: "dinner",
				MealDate: &dateStr,
			})
		}
	}
}

// Health check
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	h.respondJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}
