package api

import (
	"context"
	"net/http"

	"carro-ideal/app/internal/auth"
)

type contextKey string

const userIDKey contextKey = "userID"

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.GetUserID(r)
		if !ok {
			nextw := w
			nextw.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"success":false,"error":"unauthorized","code":"UNAUTHORIZED"}`))
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
