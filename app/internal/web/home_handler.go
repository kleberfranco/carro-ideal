package web

import (
	"net/http"

	"carro-ideal/app/internal/auth"
)

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	_, loggedIn := auth.GetUserID(r)
	render(w, "index.html", map[string]any{
		"Title":      "Carro Ideal",
		"IsLoggedIn": loggedIn,
	})
}
