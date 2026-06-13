package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Route("/auth", func(r chi.Router) {
		r.Use(JSONMiddleware)
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Post("/logout", h.Logout)
		r.Get("/me", h.Me)
	})

	r.Route("/user", func(r chi.Router) {
		r.Use(JSONMiddleware)
		r.Get("/", RequireAuth(http.HandlerFunc(h.Placeholder)))
	})
}
