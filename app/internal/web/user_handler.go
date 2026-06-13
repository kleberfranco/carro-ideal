package web

import (
	"net/http"

	"carro-ideal/app/internal/auth"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if _, ok := auth.GetUserID(r); ok {
		http.Redirect(w, r, "/recommend", http.StatusSeeOther)
		return
	}

	render(w, "login.html", map[string]any{
		"Title":      "Login",
		"IsLoggedIn": false,
	})
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if _, ok := auth.GetUserID(r); ok {
		http.Redirect(w, r, "/recommend", http.StatusSeeOther)
		return
	}

	render(w, "register.html", map[string]any{
		"Title":      "Criar conta",
		"IsLoggedIn": false,
	})
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	auth.ClearUserSession(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
