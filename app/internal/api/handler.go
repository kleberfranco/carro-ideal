package api

import (
	"net/http"

	"carro-ideal/app/internal/response"
	"carro-ideal/app/service"
)

type Handler struct {
	UserService           *service.UserService
	AuthService           *service.AuthService
	QuestionnaireService  *service.QuestionnaireService
	RecommendationService *service.RecommendationService
	VehicleService        *service.VehicleService
	SecureCookie          bool
}

func NewHandler(
	userService *service.UserService,
	authService *service.AuthService,
	questionnaireService *service.QuestionnaireService,
	recommendationService *service.RecommendationService,
	vehicleService *service.VehicleService,
	secureCookie bool,
) *Handler {
	return &Handler{
		UserService:           userService,
		AuthService:           authService,
		QuestionnaireService:  questionnaireService,
		RecommendationService: recommendationService,
		VehicleService:        vehicleService,
		SecureCookie:          secureCookie,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	userHandler := &UserHandler{
		userService:  h.UserService,
		authService:  h.AuthService,
		secureCookie: h.SecureCookie,
	}
	userHandler.Register(w, r)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	userHandler := &UserHandler{
		userService:  h.UserService,
		authService:  h.AuthService,
		secureCookie: h.SecureCookie,
	}
	userHandler.Login(w, r)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	userHandler := &UserHandler{
		userService:  h.UserService,
		authService:  h.AuthService,
		secureCookie: h.SecureCookie,
	}
	userHandler.Logout(w, r)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userHandler := &UserHandler{
		userService:  h.UserService,
		authService:  h.AuthService,
		secureCookie: h.SecureCookie,
	}
	userHandler.Me(w, r)
}

func (h *Handler) Placeholder(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotImplemented, "endpoint não implementado", "NOT_IMPLEMENTED")
}
