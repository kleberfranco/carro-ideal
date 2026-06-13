package api

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"

	"carro-ideal/app/internal/auth"
	"carro-ideal/app/internal/response"
	"carro-ideal/app/models"
	"carro-ideal/app/service"
)

type UserHandler struct {
	userService  *service.UserService
	authService  *service.AuthService
	secureCookie bool
}

type registerRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "payload inválido", "INVALID_INPUT")
		return
	}

	errors := map[string]string{}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if req.Name == "" {
		errors["name"] = "Informe seu nome completo."
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		errors["email"] = "Informe um e-mail válido."
	}
	if len(req.Password) < 8 {
		errors["password"] = "A senha deve ter pelo menos 8 caracteres."
	}
	if req.Password != req.ConfirmPassword {
		errors["confirm_password"] = "As senhas não conferem."
	}

	if len(errors) > 0 {
		response.ValidationError(w, errors)
		return
	}

	user, err := h.userService.Register(r.Context(), req.Name, req.Email, req.Password, req.ConfirmPassword)
	if err != nil {
		if service.IsEmailAlreadyUsed(err) {
			response.Error(w, http.StatusConflict, "Este e-mail já está em uso.", "EMAIL_EXISTS")
			return
		}

		response.Error(w, http.StatusInternalServerError, "Erro ao cadastrar usuário.", "INTERNAL_ERROR")
		return
	}

	token, expiresAt, err := h.authService.CreateSession(r.Context(), user.ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Usuário criado, mas houve falha ao iniciar a sessão.", "SESSION_ERROR")
		return
	}

	auth.SetSessionCookie(w, token, expiresAt, h.secureCookie)
	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"user":    userResponse(user),
		"message": "Usuário registrado e autenticado com sucesso.",
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "payload inválido", "INVALID_INPUT")
		return
	}

	errors := map[string]string{}

	if req.Email == "" {
		errors["email"] = "Informe um e-mail válido."
	}
	if req.Password == "" {
		errors["password"] = "Informe sua senha."
	}

	if len(errors) > 0 {
		response.ValidationError(w, errors)
		return
	}

	user, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error(), "AUTH_FAILED")
		return
	}

	token, expiresAt, err := h.authService.CreateSession(r.Context(), user.ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Falha ao iniciar a sessão.", "SESSION_ERROR")
		return
	}

	auth.SetSessionCookie(w, token, expiresAt, h.secureCookie)
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Login realizado com sucesso.",
		"user":    userResponse(user),
	})
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if token, ok := auth.SessionToken(r); ok {
		if err := h.authService.DestroySession(r.Context(), token); err != nil {
			response.Error(w, http.StatusInternalServerError, "Falha ao encerrar a sessão.", "SESSION_ERROR")
			return
		}
	}

	auth.ClearSessionCookie(w, h.secureCookie)
	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Logout realizado com sucesso.",
	})
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	token, ok := auth.SessionToken(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "usuário não autenticado", "UNAUTHENTICATED")
		return
	}

	userID, err := h.authService.Authenticate(r.Context(), token)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "sessão inválida ou expirada", "UNAUTHENTICATED")
		return
	}

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Falha ao buscar dados do usuário.", "INTERNAL_ERROR")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"user": userResponse(user),
	})
}

func userResponse(user *models.User) map[string]interface{} {
	return map[string]interface{}{
		"id":     user.ID,
		"name":   user.Name,
		"email":  user.Email,
		"role":   user.Role,
		"active": user.Active,
	}
}
