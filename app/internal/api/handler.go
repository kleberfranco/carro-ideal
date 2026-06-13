package api

import (
	"net/http"

	"carro-ideal/app/service"
)

type Handler struct {
	UserService *service.UserService
}

func NewHandler(userService *service.UserService) *Handler {
	return &Handler{
		UserService: userService,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	userHandler := &UserHandler{
		userService: h.UserService,
	}
	userHandler.Register(w, r)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	userHandler := &UserHandler{
		userService: h.UserService,
	}
	userHandler.Login(w, r)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	userHandler := &UserHandler{
		userService: h.UserService,
	}
	userHandler.Logout(w, r)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userHandler := &UserHandler{
		userService: h.UserService,
	}
	userHandler.Me(w, r)
}

func (h *Handler) Placeholder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	_, _ = w.Write([]byte(`{"success":false,"error":"endpoint não implementado","code":"NOT_IMPLEMENTED"}`))
}
