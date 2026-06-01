package api

import (
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, h *Handler) {
	api := r.PathPrefix("/api").Subrouter()
	api.Use(jsonMiddleware)

	api.HandleFunc("/register", h.Register).Methods("POST")
	api.HandleFunc("/login", h.Login).Methods("POST")
}
