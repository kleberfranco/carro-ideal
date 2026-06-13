package web

import (
	"net/http"

	"carro-ideal/app/internal/auth"
)

func (h *Handler) RecommendHandler(w http.ResponseWriter, r *http.Request) {
	if _, ok := auth.GetUserID(r); !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	render(w, "recommend.html", map[string]any{
		"Title":      "Recomendar Carro",
		"IsLoggedIn": true,
	})
}
