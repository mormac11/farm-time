package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"farm-time/internal/models"
)

// Todo handlers

func (h *Handler) ListTodos(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")

	todos, err := h.db.GetTodosByEvent(r.Context(), eventID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to list todos")
		return
	}
	if todos == nil {
		todos = []models.Todo{}
	}
	h.respondJSON(w, http.StatusOK, todos)
}

func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")

	var req models.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" {
		h.respondError(w, http.StatusBadRequest, "Title is required")
		return
	}

	todo, err := h.db.CreateTodo(r.Context(), eventID, req)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create todo")
		return
	}

	h.respondJSON(w, http.StatusCreated, todo)
}

func (h *Handler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "todoId")

	var req models.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	todo, err := h.db.UpdateTodo(r.Context(), todoID, req)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "Todo not found")
		return
	}

	h.respondJSON(w, http.StatusOK, todo)
}

func (h *Handler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "todoId")

	if err := h.db.DeleteTodo(r.Context(), todoID); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to delete todo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetEventWithAll returns event with attendees, meals, and todos
func (h *Handler) GetEventWithAll(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	event, err := h.db.GetEventWithAll(r.Context(), id)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "Event not found")
		return
	}

	h.respondJSON(w, http.StatusOK, event)
}
