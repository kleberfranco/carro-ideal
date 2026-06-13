package api

import (
	"encoding/json"
	"net/http"

	"carro-ideal/app/internal/auth"
	"carro-ideal/app/service"
)

type UserHandler struct {
	userService *service.UserService
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
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "payload inválido",
		})
		return
	}

	errors := map[string]string{}

	if req.Name == "" {
		errors["name"] = "Informe seu nome completo."
	}
	if req.Email == "" {
		errors["email"] = "Informe um e-mail válido."
	}
	if len(req.Password) < 6 {
		errors["password"] = "A senha deve ter pelo menos 6 caracteres."
	}
	if req.Password != req.ConfirmPassword {
		errors["confirm_password"] = "As senhas não conferem."
	}

	if len(errors) > 0 {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]map[string]string{
			"errors": errors,
		})
		return
	}

	err := h.userService.Register(r.Context(), req.Name, req.Email, req.Password, req.ConfirmPassword)
	if err != nil {
		if service.IsEmailAlreadyUsed(err) {
			writeJSON(w, http.StatusConflict, map[string]string{
				"error": "Este e-mail já está em uso.",
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Erro ao cadastrar usuário.",
		})
		return
	}

	user, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "Usuário criado, mas houve falha no login automático.",
		})
		return
	}

	auth.SetUserSession(w, user.ID)
	writeJSON(w, http.StatusCreated, map[string]string{
		"message": "Usuário registrado e autenticado com sucesso.",
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "payload inválido",
		})
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
		writeJSON(w, http.StatusUnprocessableEntity, map[string]map[string]string{
			"errors": errors,
		})
		return
	}

	user, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
		return
	}

	auth.SetUserSession(w, user.ID)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Login realizado com sucesso.",
		"user": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	auth.ClearUserSession(w)
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Logout realizado com sucesso.",
	})
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "usuário não autenticado",
		})
		return
	}

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Falha ao buscar dados do usuário.",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
			"active": user.Active,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.WriteHeader(status)
	if body == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(body)
}
