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
		`SELECT id, google_id, email, name, picture, created_at, updated_at
		 FROM users WHERE google_id = $1`, googleID,
	).Scan(&user.ID, &user.GoogleID, &user.Email, &user.Name,
		&user.Picture, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := db.pool.QueryRow(ctx,
		`SELECT id, google_id, email, name, picture, created_at, updated_at
		 FROM users WHERE id = $1`, id,
	).Scan(&user.ID, &user.GoogleID, &user.Email, &user.Name,
		&user.Picture, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) CreateOrUpdateUser(ctx context.Context, googleID, email, name, picture string) (*models.User, error) {
	now := time.Now()
	id := uuid.New().String()

	_, err := db.pool.Exec(ctx,
		`INSERT INTO users (id, google_id, email, name, picture, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $6)
		 ON CONFLICT (google_id) DO UPDATE SET
		   email = EXCLUDED.email,
		   name = EXCLUDED.name,
		   picture = EXCLUDED.picture,
		   updated_at = EXCLUDED.updated_at`,
		id, googleID, email, name, picture, now,
	)
	if err != nil {
		return nil, err
	}

	return db.GetUserByGoogleID(ctx, googleID)
}
