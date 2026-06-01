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
		"Title": "Login",
	})
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if _, ok := auth.GetUserID(r); ok {
		http.Redirect(w, r, "/recommend", http.StatusSeeOther)
		return
	}

	render(w, "register.html", map[string]any{
		"Title": "Criar conta",
	})
}
