package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"carro-ideal/app/repository"
)

const defaultSessionTTL = 24 * time.Hour

var ErrInvalidSession = errors.New("invalid or expired session")

type AuthService struct {
	repo repository.SessionRepository
	ttl  time.Duration
	now  func() time.Time
}

func NewAuthService(repo repository.SessionRepository) *AuthService {
	return &AuthService{
		repo: repo,
		ttl:  defaultSessionTTL,
		now:  time.Now,
	}
}

func (s *AuthService) CreateSession(ctx context.Context, userID int64) (string, time.Time, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", time.Time{}, fmt.Errorf("generate session token: %w", err)
	}

	token := base64.RawURLEncoding.EncodeToString(randomBytes)
	expiresAt := s.now().UTC().Add(s.ttl)
	if err := s.repo.Create(ctx, hashToken(token), userID, expiresAt); err != nil {
		return "", time.Time{}, fmt.Errorf("store session: %w", err)
	}

	return token, expiresAt, nil
}

func (s *AuthService) Authenticate(ctx context.Context, token string) (int64, error) {
	if token == "" {
		return 0, ErrInvalidSession
	}

	userID, err := s.repo.GetUserID(ctx, hashToken(token), s.now().UTC())
	if err != nil {
		if errors.Is(err, repository.ErrSessionNotFound) {
			return 0, ErrInvalidSession
		}
		return 0, fmt.Errorf("load session: %w", err)
	}
	return userID, nil
}

func (s *AuthService) DestroySession(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	return s.repo.Delete(ctx, hashToken(token))
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
