package admin

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, h *Handler) {
	s := r.PathPrefix("/admin").Subrouter()
	s.HandleFunc("", h.AdminHandler).Methods(http.MethodGet)
	s.HandleFunc("/", h.AdminHandler).Methods(http.MethodGet)
}
