package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrSessionNotFound = errors.New("session not found")

type SessionRepository interface {
	Create(ctx context.Context, tokenHash string, userID int64, expiresAt time.Time) error
	GetUserID(ctx context.Context, tokenHash string, now time.Time) (int64, error)
	Delete(ctx context.Context, tokenHash string) error
	DeleteExpired(ctx context.Context, now time.Time) error
}

type sessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, tokenHash string, userID int64, expiresAt time.Time) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO sessions (token_hash, user_id, expires_at) VALUES ($1, $2, $3)",
		tokenHash, userID, expiresAt,
	)
	return err
}

func (r *sessionRepository) GetUserID(ctx context.Context, tokenHash string, now time.Time) (int64, error) {
	var userID int64
	err := r.db.QueryRowContext(
		ctx,
		`SELECT s.user_id
		 FROM sessions s
		 JOIN users u ON u.id = s.user_id
		 WHERE s.token_hash=$1 AND s.expires_at>$2 AND u.active=true`,
		tokenHash, now,
	).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrSessionNotFound
	}
	return userID, err
}

func (r *sessionRepository) Delete(ctx context.Context, tokenHash string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM sessions WHERE token_hash=$1", tokenHash)
	return err
}

func (r *sessionRepository) DeleteExpired(ctx context.Context, now time.Time) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM sessions WHERE expires_at<=$1", now)
	return err
}
