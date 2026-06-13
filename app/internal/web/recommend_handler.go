package web

import (
	"net/http"
)

func (h *Handler) RecommendHandler(w http.ResponseWriter, r *http.Request) {
	if !h.authenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	render(w, "recommend.html", map[string]any{
		"Title":      "Recomendar Carro",
		"IsLoggedIn": true,
	})
}
