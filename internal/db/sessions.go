package db

import (
	"context"
	"time"

	"farm-time/internal/models"
)

func (db *DB) CreateSession(ctx context.Context, sessionID, userID string, expiresAt time.Time) error {
	_, err := db.pool.Exec(ctx,
		`INSERT INTO sessions (id, user_id, expires_at, created_at)
		 VALUES ($1, $2, $3, NOW())`,
		sessionID, userID, expiresAt,
	)
	return err
}

func (db *DB) GetSession(ctx context.Context, sessionID string) (*models.Session, error) {
	var session models.Session
	err := db.pool.QueryRow(ctx,
		`SELECT id, user_id, expires_at, created_at
		 FROM sessions WHERE id = $1`, sessionID,
	).Scan(&session.ID, &session.UserID, &session.ExpiresAt, &session.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (db *DB) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := db.pool.Exec(ctx, `DELETE FROM sessions WHERE id = $1`, sessionID)
	return err
}

func (db *DB) DeleteUserSessions(ctx context.Context, userID string) error {
	_, err := db.pool.Exec(ctx, `DELETE FROM sessions WHERE user_id = $1`, userID)
	return err
}

func (db *DB) CleanExpiredSessions(ctx context.Context) error {
	_, err := db.pool.Exec(ctx, `DELETE FROM sessions WHERE expires_at < NOW()`)
	return err
}
