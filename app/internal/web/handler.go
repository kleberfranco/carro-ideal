package web

import (
	"net/http"

	"carro-ideal/app/internal/auth"
	"carro-ideal/app/service"
)

type Handler struct {
	UserService  *service.UserService
	AuthService  *service.AuthService
	SecureCookie bool
}

func NewHandler(userService *service.UserService, authService *service.AuthService, secureCookie bool) *Handler {
	return &Handler{
		UserService:  userService,
		AuthService:  authService,
		SecureCookie: secureCookie,
	}
}

func (h *Handler) authenticated(r *http.Request) bool {
	token, ok := auth.SessionToken(r)
	if !ok {
		return false
	}

	_, err := h.AuthService.Authenticate(r.Context(), token)
	return err == nil
}
