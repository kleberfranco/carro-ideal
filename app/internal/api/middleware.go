package api

import (
	"context"
	"net/http"

	"carro-ideal/app/internal/auth"
	"carro-ideal/app/internal/response"
	"carro-ideal/app/service"
)

type contextKey string

const userIDKey contextKey = "userID"

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func RequireAuth(authService *service.AuthService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := auth.SessionToken(r)
		if !ok {
			response.Error(w, http.StatusUnauthorized, "autenticação necessária", "UNAUTHENTICATED")
			return
		}

		userID, err := authService.Authenticate(r.Context(), token)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "sessão inválida ou expirada", "UNAUTHENTICATED")
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(r *http.Request) (int64, bool) {
	value := r.Context().Value(userIDKey)
	if value == nil {
		return 0, false
	}

	userID, ok := value.(int64)
	return userID, ok
}
