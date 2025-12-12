package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"farm-time/internal/models"
)

type DB struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context) (*DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/farmtime?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{pool: pool}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) Migrate(ctx context.Context) error {
	schema := `
	CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		location TEXT,
		start_time TIMESTAMPTZ NOT NULL,
		end_time TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS attendees (
		id TEXT PRIMARY KEY,
		event_id TEXT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'attending',
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_attendees_event_id ON attendees(event_id);
	CREATE INDEX IF NOT EXISTS idx_events_start_time ON events(start_time);

	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		google_id TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		picture TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		expires_at TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
	CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
	CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id);

	CREATE TABLE IF NOT EXISTS meals (
		id TEXT PRIMARY KEY,
		event_id TEXT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
		name TEXT NOT NULL,
		meal_type TEXT NOT NULL,
		meal_date DATE,
		notes TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS meal_items (
		id TEXT PRIMARY KEY,
		meal_id TEXT NOT NULL REFERENCES meals(id) ON DELETE CASCADE,
		name TEXT NOT NULL,
		description TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS meal_signups (
		id TEXT PRIMARY KEY,
		meal_item_id TEXT NOT NULL REFERENCES meal_items(id) ON DELETE CASCADE,
		user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		notes TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		UNIQUE(meal_item_id, user_id)
	);

	CREATE INDEX IF NOT EXISTS idx_meals_event_id ON meals(event_id);
	CREATE INDEX IF NOT EXISTS idx_meal_items_meal_id ON meal_items(meal_id);
	CREATE INDEX IF NOT EXISTS idx_meal_signups_item ON meal_signups(meal_item_id);
	CREATE INDEX IF NOT EXISTS idx_meal_signups_user ON meal_signups(user_id);
	`

	_, err := db.pool.Exec(ctx, schema)
	return err
}

// Event operations

func (db *DB) CreateEvent(ctx context.Context, req models.CreateEventRequest) (*models.Event, error) {
	event := &models.Event{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := db.pool.Exec(ctx,
		`INSERT INTO events (id, title, description, location, start_time, end_time, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		event.ID, event.Title, event.Description, event.Location,
		event.StartTime, event.EndTime, event.CreatedAt, event.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (db *DB) GetEvent(ctx context.Context, id string) (*models.Event, error) {
	var event models.Event
	err := db.pool.QueryRow(ctx,
		`SELECT id, title, description, location, start_time, end_time, created_at, updated_at
		 FROM events WHERE id = $1`, id,
	).Scan(&event.ID, &event.Title, &event.Description, &event.Location,
		&event.StartTime, &event.EndTime, &event.CreatedAt, &event.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (db *DB) ListEvents(ctx context.Context) ([]models.Event, error) {
	rows, err := db.pool.Query(ctx,
		`SELECT id, title, description, location, start_time, end_time, created_at, updated_at
		 FROM events ORDER BY start_time ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		if err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.Location,
			&e.StartTime, &e.EndTime, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (db *DB) UpdateEvent(ctx context.Context, id string, req models.UpdateEventRequest) (*models.Event, error) {
	event, err := db.GetEvent(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		event.Title = *req.Title
	}
	if req.Description != nil {
		event.Description = *req.Description
	}
	if req.Location != nil {
		event.Location = *req.Location
	}
	if req.StartTime != nil {
		event.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		event.EndTime = *req.EndTime
	}
	event.UpdatedAt = time.Now()

	_, err = db.pool.Exec(ctx,
		`UPDATE events SET title=$1, description=$2, location=$3, start_time=$4, end_time=$5, updated_at=$6
		 WHERE id=$7`,
		event.Title, event.Description, event.Location, event.StartTime, event.EndTime, event.UpdatedAt, id,
	)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (db *DB) DeleteEvent(ctx context.Context, id string) error {
	_, err := db.pool.Exec(ctx, `DELETE FROM events WHERE id = $1`, id)
	return err
}

// Attendee operations

func (db *DB) CreateAttendee(ctx context.Context, eventID string, req models.CreateAttendeeRequest) (*models.Attendee, error) {
	attendee := &models.Attendee{
		ID:        uuid.New().String(),
		EventID:   eventID,
		Name:      req.Name,
		Email:     req.Email,
		Status:    req.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if attendee.Status == "" {
		attendee.Status = "attending"
	}

	_, err := db.pool.Exec(ctx,
		`INSERT INTO attendees (id, event_id, name, email, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		attendee.ID, attendee.EventID, attendee.Name, attendee.Email,
		attendee.Status, attendee.CreatedAt, attendee.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return attendee, nil
}

func (db *DB) GetAttendeesByEvent(ctx context.Context, eventID string) ([]models.Attendee, error) {
	rows, err := db.pool.Query(ctx,
		`SELECT id, event_id, name, email, status, created_at, updated_at
		 FROM attendees WHERE event_id = $1 ORDER BY name ASC`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendees []models.Attendee
	for rows.Next() {
		var a models.Attendee
		if err := rows.Scan(&a.ID, &a.EventID, &a.Name, &a.Email,
			&a.Status, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		attendees = append(attendees, a)
	}

	return attendees, nil
}

func (db *DB) UpdateAttendee(ctx context.Context, id string, req models.UpdateAttendeeRequest) (*models.Attendee, error) {
	var attendee models.Attendee
	err := db.pool.QueryRow(ctx,
		`SELECT id, event_id, name, email, status, created_at, updated_at
		 FROM attendees WHERE id = $1`, id,
	).Scan(&attendee.ID, &attendee.EventID, &attendee.Name, &attendee.Email,
		&attendee.Status, &attendee.CreatedAt, &attendee.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		attendee.Name = *req.Name
	}
	if req.Email != nil {
		attendee.Email = *req.Email
	}
	if req.Status != nil {
		attendee.Status = *req.Status
	}
	attendee.UpdatedAt = time.Now()

	_, err = db.pool.Exec(ctx,
		`UPDATE attendees SET name=$1, email=$2, status=$3, updated_at=$4 WHERE id=$5`,
		attendee.Name, attendee.Email, attendee.Status, attendee.UpdatedAt, id,
	)
	if err != nil {
		return nil, err
	}

	return &attendee, nil
}

func (db *DB) DeleteAttendee(ctx context.Context, id string) error {
	_, err := db.pool.Exec(ctx, `DELETE FROM attendees WHERE id = $1`, id)
	return err
}

func (db *DB) GetEventWithAttendees(ctx context.Context, id string) (*models.EventWithAttendees, error) {
	event, err := db.GetEvent(ctx, id)
	if err != nil {
		return nil, err
	}

	attendees, err := db.GetAttendeesByEvent(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.EventWithAttendees{
		Event:     *event,
		Attendees: attendees,
	}, nil
}
