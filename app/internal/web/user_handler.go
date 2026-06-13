package web

import (
	"net/http"

	"carro-ideal/app/internal/auth"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if h.authenticated(r) {
		http.Redirect(w, r, "/recommend", http.StatusSeeOther)
		return
	}

	render(w, "login.html", map[string]any{
		"Title":      "Login",
		"IsLoggedIn": false,
	})
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if h.authenticated(r) {
		http.Redirect(w, r, "/recommend", http.StatusSeeOther)
		return
	}

	render(w, "register.html", map[string]any{
		"Title":      "Criar conta",
		"IsLoggedIn": false,
	})
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if token, ok := auth.SessionToken(r); ok {
		_ = h.AuthService.DestroySession(r.Context(), token)
	}
	auth.ClearSessionCookie(w, h.SecureCookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
