package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"carro-ideal/app/internal/auth"
	"carro-ideal/app/models"
	"carro-ideal/app/repository"
	"carro-ideal/app/service"
)

type middlewareSessionRepository struct {
	tokenHash string
	userID    int64
	expiresAt time.Time
}

type roleUserRepository struct {
	user *models.User
}

func (r *roleUserRepository) ExistsByEmail(context.Context, string) (bool, error) {
	return false, nil
}
func (r *roleUserRepository) Create(context.Context, *models.User) error { return nil }
func (r *roleUserRepository) GetByEmail(context.Context, string) (*models.User, error) {
	return r.user, nil
}
func (r *roleUserRepository) GetByID(context.Context, int64) (*models.User, error) {
	return r.user, nil
}
func (r *roleUserRepository) Update(context.Context, *models.User) error { return nil }
func (r *roleUserRepository) Deactivate(context.Context, int64) error    { return nil }

func (r *middlewareSessionRepository) Create(_ context.Context, tokenHash string, userID int64, expiresAt time.Time) error {
	r.tokenHash = tokenHash
	r.userID = userID
	r.expiresAt = expiresAt
	return nil
}

func (r *middlewareSessionRepository) GetUserID(_ context.Context, tokenHash string, now time.Time) (int64, error) {
	if tokenHash != r.tokenHash || !now.Before(r.expiresAt) {
		return 0, repository.ErrSessionNotFound
	}
	return r.userID, nil
}

func (r *middlewareSessionRepository) Delete(_ context.Context, _ string) error {
	return nil
}

func (r *middlewareSessionRepository) DeleteExpired(_ context.Context, _ time.Time) error {
	return nil
}

func TestRequireAuth(t *testing.T) {
	repo := &middlewareSessionRepository{}
	authService := service.NewAuthService(repo)
	token, expiresAt, err := authService.CreateSession(context.Background(), 7)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	protected := RequireAuth(authService, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserIDFromContext(r)
		if !ok || userID != 7 {
			t.Fatalf("context user ID = %d, %v; want 7, true", userID, ok)
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	t.Run("valid session", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/api/user", nil)
		request.AddCookie(&http.Cookie{Name: "carro_session", Value: token, Expires: expiresAt})
		recorder := httptest.NewRecorder()

		protected.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
		}
	})

	t.Run("missing session", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/api/user", nil)
		recorder := httptest.NewRecorder()

		protected.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
		}
	})

	t.Run("expired session", func(t *testing.T) {
		repo.expiresAt = time.Now().Add(-time.Minute)
		request := httptest.NewRequest(http.MethodGet, "/api/user", nil)
		request.AddCookie(&http.Cookie{Name: "carro_session", Value: token})
		recorder := httptest.NewRecorder()

		protected.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
		}
	})
}

func TestSessionCookieSecurityFlags(t *testing.T) {
	recorder := httptest.NewRecorder()
	auth.SetSessionCookie(recorder, "token", time.Now().Add(time.Hour), true)

	cookies := recorder.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("cookie count = %d, want 1", len(cookies))
	}
	if !cookies[0].HttpOnly || !cookies[0].Secure || cookies[0].SameSite != http.SameSiteLaxMode {
		t.Fatal("session cookie is missing required security flags")
	}
}

func TestRequireAdmin(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	tests := []struct {
		name       string
		role       string
		wantStatus int
	}{
		{name: "admin accepted", role: "admin", wantStatus: http.StatusNoContent},
		{name: "regular user rejected", role: "user", wantStatus: http.StatusForbidden},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userService := service.NewUserService(&roleUserRepository{
				user: &models.User{ID: 7, Role: test.role, Active: true},
			})
			handler := RequireAdmin(userService, next)
			request := httptest.NewRequest(http.MethodGet, "/api/admin/dashboard", nil)
			request = request.WithContext(context.WithValue(request.Context(), userIDKey, int64(7)))
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)

			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, test.wantStatus)
			}
		})
	}
}
