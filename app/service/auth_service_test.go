package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"carro-ideal/app/repository"
)

type fakeSessionRepository struct {
	tokenHash string
	userID    int64
	expiresAt time.Time
	deleted   string
	getErr    error
}

func (r *fakeSessionRepository) Create(_ context.Context, tokenHash string, userID int64, expiresAt time.Time) error {
	r.tokenHash = tokenHash
	r.userID = userID
	r.expiresAt = expiresAt
	return nil
}

func (r *fakeSessionRepository) GetUserID(_ context.Context, tokenHash string, now time.Time) (int64, error) {
	if r.getErr != nil {
		return 0, r.getErr
	}
	if tokenHash != r.tokenHash || !now.Before(r.expiresAt) {
		return 0, repository.ErrSessionNotFound
	}
	return r.userID, nil
}

func (r *fakeSessionRepository) Delete(_ context.Context, tokenHash string) error {
	r.deleted = tokenHash
	return nil
}

func (r *fakeSessionRepository) DeleteExpired(_ context.Context, _ time.Time) error {
	return nil
}

func TestAuthServiceSessionLifecycle(t *testing.T) {
	repo := &fakeSessionRepository{}
	service := NewAuthService(repo)
	now := time.Date(2026, 6, 13, 12, 0, 0, 0, time.UTC)
	service.now = func() time.Time { return now }

	token, expiresAt, err := service.CreateSession(context.Background(), 42)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}
	if token == "" {
		t.Fatal("CreateSession() returned an empty token")
	}
	if repo.tokenHash == token {
		t.Fatal("repository stored the raw session token")
	}
	if want := now.Add(24 * time.Hour); !expiresAt.Equal(want) {
		t.Fatalf("expiresAt = %v, want %v", expiresAt, want)
	}

	userID, err := service.Authenticate(context.Background(), token)
	if err != nil {
		t.Fatalf("Authenticate() error = %v", err)
	}
	if userID != 42 {
		t.Fatalf("Authenticate() userID = %d, want 42", userID)
	}

	if err := service.DestroySession(context.Background(), token); err != nil {
		t.Fatalf("DestroySession() error = %v", err)
	}
	if repo.deleted != repo.tokenHash {
		t.Fatal("DestroySession() did not delete the hashed token")
	}
}

func TestAuthServiceRejectsExpiredSession(t *testing.T) {
	repo := &fakeSessionRepository{getErr: repository.ErrSessionNotFound}
	service := NewAuthService(repo)

	_, err := service.Authenticate(context.Background(), "expired-token")
	if !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("Authenticate() error = %v, want ErrInvalidSession", err)
	}
}
