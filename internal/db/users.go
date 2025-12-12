package db

import (
	"context"
	"time"

	"github.com/google/uuid"

	"farm-time/internal/models"
)

func (db *DB) GetUserByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	var user models.User
	err := db.pool.QueryRow(ctx,
		`SELECT id, google_id, email, name, picture, is_admin, can_create_events, created_at, updated_at
		 FROM users WHERE google_id = $1`, googleID,
	).Scan(&user.ID, &user.GoogleID, &user.Email, &user.Name,
		&user.Picture, &user.IsAdmin, &user.CanCreateEvents, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := db.pool.QueryRow(ctx,
		`SELECT id, google_id, email, name, picture, is_admin, can_create_events, created_at, updated_at
		 FROM users WHERE id = $1`, id,
	).Scan(&user.ID, &user.GoogleID, &user.Email, &user.Name,
		&user.Picture, &user.IsAdmin, &user.CanCreateEvents, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) CreateOrUpdateUser(ctx context.Context, googleID, email, name, picture string) (*models.User, error) {
	now := time.Now()
	id := uuid.New().String()

	// Check if this is the first user - make them admin
	var userCount int
	db.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&userCount)
	isFirstUser := userCount == 0

	_, err := db.pool.Exec(ctx,
		`INSERT INTO users (id, google_id, email, name, picture, is_admin, can_create_events, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
		 ON CONFLICT (google_id) DO UPDATE SET
		   email = EXCLUDED.email,
		   name = EXCLUDED.name,
		   picture = EXCLUDED.picture,
		   updated_at = EXCLUDED.updated_at`,
		id, googleID, email, name, picture, isFirstUser, isFirstUser, now,
	)
	if err != nil {
		return nil, err
	}

	return db.GetUserByGoogleID(ctx, googleID)
}

// Admin operations

func (db *DB) ListAllUsers(ctx context.Context) ([]models.User, error) {
	rows, err := db.pool.Query(ctx,
		`SELECT id, google_id, email, name, picture, is_admin, can_create_events, created_at, updated_at
		 FROM users ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.GoogleID, &u.Email, &u.Name,
			&u.Picture, &u.IsAdmin, &u.CanCreateEvents, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (db *DB) UpdateUserPermissions(ctx context.Context, userID string, req models.UpdateUserPermissionsRequest) (*models.User, error) {
	user, err := db.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.CanCreateEvents != nil {
		user.CanCreateEvents = *req.CanCreateEvents
	}
	user.UpdatedAt = time.Now()

	_, err = db.pool.Exec(ctx,
		`UPDATE users SET can_create_events=$1, updated_at=$2 WHERE id=$3`,
		user.CanCreateEvents, user.UpdatedAt, userID,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
