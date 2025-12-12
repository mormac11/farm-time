package db

import (
	"context"
	"time"

	"github.com/google/uuid"

	"farm-time/internal/models"
)

// Todo operations

func (db *DB) CreateTodo(ctx context.Context, eventID string, req models.CreateTodoRequest) (*models.Todo, error) {
	todo := &models.Todo{
		ID:                 uuid.New().String(),
		EventID:            eventID,
		Title:              req.Title,
		Description:        req.Description,
		Completed:          false,
		AssignedAttendeeID: req.AssignedAttendeeID,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	_, err := db.pool.Exec(ctx,
		`INSERT INTO todos (id, event_id, title, description, completed, assigned_attendee_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		todo.ID, todo.EventID, todo.Title, todo.Description, todo.Completed, todo.AssignedAttendeeID, todo.CreatedAt, todo.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get attendee name if assigned
	if todo.AssignedAttendeeID != nil {
		var name string
		err := db.pool.QueryRow(ctx, `SELECT name FROM attendees WHERE id = $1`, *todo.AssignedAttendeeID).Scan(&name)
		if err == nil {
			todo.AssignedAttendeeName = &name
		}
	}

	return todo, nil
}

func (db *DB) GetTodo(ctx context.Context, id string) (*models.Todo, error) {
	var todo models.Todo
	err := db.pool.QueryRow(ctx,
		`SELECT t.id, t.event_id, t.title, COALESCE(t.description, ''), t.completed, t.assigned_attendee_id, a.name, t.created_at, t.updated_at
		 FROM todos t
		 LEFT JOIN attendees a ON t.assigned_attendee_id = a.id
		 WHERE t.id = $1`, id,
	).Scan(&todo.ID, &todo.EventID, &todo.Title, &todo.Description, &todo.Completed, &todo.AssignedAttendeeID, &todo.AssignedAttendeeName, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (db *DB) GetTodosByEvent(ctx context.Context, eventID string) ([]models.Todo, error) {
	rows, err := db.pool.Query(ctx,
		`SELECT t.id, t.event_id, t.title, COALESCE(t.description, ''), t.completed, t.assigned_attendee_id, a.name, t.created_at, t.updated_at
		 FROM todos t
		 LEFT JOIN attendees a ON t.assigned_attendee_id = a.id
		 WHERE t.event_id = $1 ORDER BY t.completed ASC, t.created_at ASC`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.ID, &t.EventID, &t.Title, &t.Description, &t.Completed, &t.AssignedAttendeeID, &t.AssignedAttendeeName, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	return todos, nil
}

func (db *DB) UpdateTodo(ctx context.Context, id string, req models.UpdateTodoRequest) (*models.Todo, error) {
	todo, err := db.GetTodo(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}
	if req.AssignedAttendeeID != nil {
		if *req.AssignedAttendeeID == "" {
			todo.AssignedAttendeeID = nil
			todo.AssignedAttendeeName = nil
		} else {
			todo.AssignedAttendeeID = req.AssignedAttendeeID
			// Get attendee name
			var name string
			err := db.pool.QueryRow(ctx, `SELECT name FROM attendees WHERE id = $1`, *req.AssignedAttendeeID).Scan(&name)
			if err == nil {
				todo.AssignedAttendeeName = &name
			}
		}
	}
	todo.UpdatedAt = time.Now()

	_, err = db.pool.Exec(ctx,
		`UPDATE todos SET title=$1, description=$2, completed=$3, assigned_attendee_id=$4, updated_at=$5 WHERE id=$6`,
		todo.Title, todo.Description, todo.Completed, todo.AssignedAttendeeID, todo.UpdatedAt, id,
	)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (db *DB) DeleteTodo(ctx context.Context, id string) error {
	_, err := db.pool.Exec(ctx, `DELETE FROM todos WHERE id = $1`, id)
	return err
}

// GetEventWithAll returns event with attendees, meals, and todos
func (db *DB) GetEventWithAll(ctx context.Context, id string) (*models.EventWithAll, error) {
	eventWithMeals, err := db.GetEventWithMeals(ctx, id)
	if err != nil {
		return nil, err
	}

	todos, err := db.GetTodosByEvent(ctx, id)
	if err != nil {
		return nil, err
	}
	if todos == nil {
		todos = []models.Todo{}
	}

	return &models.EventWithAll{
		EventWithMeals: *eventWithMeals,
		Todos:          todos,
	}, nil
}
