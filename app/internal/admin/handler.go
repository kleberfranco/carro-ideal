package admin

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"carro-ideal/app/internal/auth"
	"carro-ideal/app/internal/response"
	"carro-ideal/app/models"
	"carro-ideal/app/repository"
	"carro-ideal/app/service"

	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
)

type Handler struct {
	UserService  *service.UserService
	AuthService  *service.AuthService
	AdminService *service.AdminService
}

func NewHandler(
	userService *service.UserService,
	authService *service.AuthService,
	adminService *service.AdminService,
) *Handler {
	return &Handler{
		UserService:  userService,
		AuthService:  authService,
		AdminService: adminService,
	}
}

func (h *Handler) Page(w http.ResponseWriter, r *http.Request) {
	token, ok := auth.SessionToken(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	userID, err := h.AuthService.Authenticate(r.Context(), token)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	user, err := h.UserService.GetByID(r.Context(), userID)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !strings.EqualFold(user.Role, "admin") {
		http.Error(w, "acesso restrito a administradores", http.StatusForbidden)
		return
	}

	t, err := template.ParseFiles("web/templates/layout.html", "web/templates/admin.html")
	if err != nil {
		http.Error(w, "falha ao carregar painel", http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, map[string]interface{}{
		"Title":      "Administração",
		"IsLoggedIn": true,
		"IsAdmin":    true,
	}); err != nil {
		http.Error(w, "falha ao renderizar painel", http.StatusInternalServerError)
	}
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	stats, err := h.AdminService.Stats(r.Context())
	if err != nil {
		adminError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, stats)
}

func (h *Handler) Vehicles(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 1)
	limit := queryInt(r, "limit", 10)
	items, total, err := h.AdminService.Vehicles(r.Context(), r.URL.Query().Get("search"), page, limit)
	if err != nil {
		adminError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *Handler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	var input models.VehicleInput
	if !decode(w, r, &input) {
		return
	}
	item, err := h.AdminService.CreateVehicle(r.Context(), input)
	if err != nil {
		adminError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	var input models.VehicleInput
	if !decode(w, r, &input) {
		return
	}
	item, err := h.AdminService.UpdateVehicle(r.Context(), id, input)
	if err != nil {
		adminError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	if err := h.AdminService.DeleteVehicle(r.Context(), id); err != nil {
		adminError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{"id": id, "active": false})
}

func (h *Handler) Categories(w http.ResponseWriter, r *http.Request) {
	items, err := h.AdminService.Categories(r.Context())
	if err != nil {
		adminError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{"items": items})
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var input models.CategoryInput
	if !decode(w, r, &input) {
		return
	}
	item, err := h.AdminService.CreateCategory(r.Context(), input)
	if err != nil {
		adminError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	var input models.CategoryInput
	if !decode(w, r, &input) {
		return
	}
	item, err := h.AdminService.UpdateCategory(r.Context(), id, input)
	if err != nil {
		adminError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	if err := h.AdminService.DeleteCategory(r.Context(), id); err != nil {
		adminError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{"id": id, "active": false})
}

func (h *Handler) Questions(w http.ResponseWriter, r *http.Request) {
	items, err := h.AdminService.Questions(r.Context())
	if err != nil {
		adminError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{"items": items})
}

func (h *Handler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var input models.QuestionInput
	if !decode(w, r, &input) {
		return
	}
	item, err := h.AdminService.CreateQuestion(r.Context(), input)
	if err != nil {
		adminError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	var input models.QuestionInput
	if !decode(w, r, &input) {
		return
	}
	item, err := h.AdminService.UpdateQuestion(r.Context(), id, input)
	if err != nil {
		adminError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	if err := h.AdminService.DeleteQuestion(r.Context(), id); err != nil {
		adminError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{"id": id, "active": false})
}

func (h *Handler) CreateOption(w http.ResponseWriter, r *http.Request) {
	questionID, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	var input models.AnswerOptionInput
	if !decode(w, r, &input) {
		return
	}
	item, err := h.AdminService.CreateOption(r.Context(), questionID, input)
	if err != nil {
		adminError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateOption(w http.ResponseWriter, r *http.Request) {
	questionID, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	optionID, ok := pathID(w, r, "optionId")
	if !ok {
		return
	}
	var input models.AnswerOptionInput
	if !decode(w, r, &input) {
		return
	}
	item, err := h.AdminService.UpdateOption(r.Context(), questionID, optionID, input)
	if err != nil {
		adminError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteOption(w http.ResponseWriter, r *http.Request) {
	questionID, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	optionID, ok := pathID(w, r, "optionId")
	if !ok {
		return
	}
	if err := h.AdminService.DeleteOption(r.Context(), questionID, optionID); err != nil {
		adminError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{"id": optionID, "active": false})
}

func decode(w http.ResponseWriter, r *http.Request, target interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		response.Error(w, http.StatusBadRequest, "payload inválido", "INVALID_INPUT")
		return false
	}
	return true
}

func pathID(w http.ResponseWriter, r *http.Request, name string) (int64, bool) {
	id, err := strconv.ParseInt(chi.URLParam(r, name), 10, 64)
	if err != nil || id < 1 {
		response.Error(w, http.StatusBadRequest, "identificador inválido", "INVALID_INPUT")
		return 0, false
	}
	return id, true
}

func queryInt(r *http.Request, name string, fallback int) int {
	value, err := strconv.Atoi(r.URL.Query().Get(name))
	if err != nil || value < 1 {
		return fallback
	}
	return value
}

func adminError(w http.ResponseWriter, err error) {
	var pqErr *pq.Error
	switch {
	case errors.Is(err, service.ErrAdminValidation):
		response.Error(w, http.StatusUnprocessableEntity, err.Error(), "VALIDATION_ERROR")
	case errors.Is(err, repository.ErrVehicleNotFound),
		errors.Is(err, repository.ErrCategoryNotFound),
		errors.Is(err, repository.ErrQuestionNotFound),
		errors.Is(err, repository.ErrOptionNotFound):
		response.Error(w, http.StatusNotFound, "registro não encontrado", "NOT_FOUND")
	case errors.Is(err, repository.ErrCategoryInUse):
		response.Error(w, http.StatusConflict, "categoria possui veículos ativos", "CATEGORY_IN_USE")
	case errors.As(err, &pqErr) && pqErr.Code == "23505":
		response.Error(w, http.StatusConflict, "registro duplicado", "CONFLICT")
	case errors.As(err, &pqErr) && pqErr.Code == "23503":
		response.Error(w, http.StatusUnprocessableEntity, "relacionamento inválido", "INVALID_REFERENCE")
	default:
		response.Error(w, http.StatusInternalServerError, "falha na operação administrativa", "INTERNAL_ERROR")
	}
}
