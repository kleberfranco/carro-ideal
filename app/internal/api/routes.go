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
		r.Handle("/", RequireAuth(h.AuthService, http.HandlerFunc(h.Placeholder)))
	})

	r.With(JSONMiddleware, func(next http.Handler) http.Handler {
		return RequireAuth(h.AuthService, next)
	}).Get("/questions", h.Questions)

	r.With(JSONMiddleware, func(next http.Handler) http.Handler {
		return RequireAuth(h.AuthService, next)
	}).Get("/recommendations", h.RecommendationHistory)

	r.Route("/recommendations", func(r chi.Router) {
		r.Use(JSONMiddleware)
		r.Use(func(next http.Handler) http.Handler {
			return RequireAuth(h.AuthService, next)
		})
		r.Post("/generate", h.GenerateRecommendations)
		r.Get("/", h.RecommendationHistory)
		r.Get("/{id}", h.RecommendationDetails)
	})

	r.With(JSONMiddleware).Get("/vehicles", h.Vehicles)

	r.Route("/vehicles", func(r chi.Router) {
		r.Use(JSONMiddleware)
		r.Get("/", h.Vehicles)
		r.Get("/{id}", h.VehicleDetails)
	})
}
