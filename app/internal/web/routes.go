package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, h *Handler) {
	static := http.StripPrefix("/static/", http.FileServer(http.Dir("web/static")))
	r.PathPrefix("/static/").Handler(static)

	r.HandleFunc("/", h.HomeHandler).Methods(http.MethodGet)
	r.HandleFunc("/login", h.LoginHandler).Methods(http.MethodGet)
	r.HandleFunc("/register", h.RegisterHandler).Methods(http.MethodGet)
	r.HandleFunc("/recommend", h.RecommendHandler).Methods(http.MethodGet)
}
