package api

import (
	"context"
	"net/http"
	"strings"

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

func RequireAdmin(userService *service.UserService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserIDFromContext(r)
		if !ok {
			response.Error(w, http.StatusUnauthorized, "autenticação necessária", "UNAUTHENTICATED")
			return
		}
		user, err := userService.GetByID(r.Context(), userID)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "usuário não encontrado", "UNAUTHENTICATED")
			return
		}
		if !strings.EqualFold(user.Role, "admin") {
			response.Error(w, http.StatusForbidden, "acesso restrito a administradores", "FORBIDDEN")
			return
		}
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
