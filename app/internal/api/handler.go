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
