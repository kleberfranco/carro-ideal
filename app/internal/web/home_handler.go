package web

import (
	"net/http"
)

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "index.html", map[string]any{
		"Title":      "Carro Ideal",
		"IsLoggedIn": h.authenticated(r),
	})
}
