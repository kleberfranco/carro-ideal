package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, h *Handler) {
	static := http.StripPrefix("/static/", http.FileServer(http.Dir("web/static")))
	r.Handle("/static/*", static)

	r.Get("/", h.HomeHandler)
	r.Get("/login", h.LoginHandler)
	r.Get("/register", h.RegisterHandler)
	r.Get("/recommend", h.RecommendHandler)
}
