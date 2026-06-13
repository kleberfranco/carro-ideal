package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"carro-ideal/app/internal/response"
	"carro-ideal/app/models"
	"carro-ideal/app/repository"
	"carro-ideal/app/service"

	"github.com/go-chi/chi/v5"
)

type generateRecommendationRequest struct {
	Answers []models.SubmittedAnswer `json:"answers"`
}

func (h *Handler) Questions(w http.ResponseWriter, r *http.Request) {
	questions, err := h.QuestionnaireService.GetActive(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Falha ao carregar o questionário.", "INTERNAL_ERROR")
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{"items": questions})
}

func (h *Handler) GenerateRecommendations(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "autenticação necessária", "UNAUTHENTICATED")
		return
	}

	var request generateRecommendationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, http.StatusBadRequest, "payload inválido", "INVALID_INPUT")
		return
	}

	recommendation, err := h.RecommendationService.Generate(r.Context(), userID, request.Answers)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrIncompleteQuestionnaire):
			response.Error(w, http.StatusUnprocessableEntity, err.Error(), "INCOMPLETE_QUESTIONNAIRE")
		case errors.Is(err, service.ErrInvalidAnswer):
			response.Error(w, http.StatusUnprocessableEntity, err.Error(), "INVALID_ANSWER")
		case errors.Is(err, service.ErrNoVehicles):
			response.Error(w, http.StatusUnprocessableEntity, err.Error(), "NO_VEHICLES")
		default:
			response.Error(w, http.StatusInternalServerError, "Falha ao gerar recomendações.", "INTERNAL_ERROR")
		}
		return
	}
	response.JSON(w, http.StatusCreated, recommendation)
}

func (h *Handler) RecommendationHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "autenticação necessária", "UNAUTHENTICATED")
		return
	}
	page := queryInt(r, "page", 1)
	limit := queryInt(r, "limit", 10)
	items, total, err := h.RecommendationService.History(r.Context(), userID, page, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Falha ao carregar o histórico.", "INTERNAL_ERROR")
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *Handler) RecommendationDetails(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "autenticação necessária", "UNAUTHENTICATED")
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "identificador inválido", "INVALID_INPUT")
		return
	}
	recommendation, err := h.RecommendationService.Get(r.Context(), id, userID)
	if errors.Is(err, repository.ErrRecommendationNotFound) {
		response.Error(w, http.StatusNotFound, "Recomendação não encontrada.", "NOT_FOUND")
		return
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Falha ao carregar a recomendação.", "INTERNAL_ERROR")
		return
	}
	response.JSON(w, http.StatusOK, recommendation)
}

func (h *Handler) Vehicles(w http.ResponseWriter, r *http.Request) {
	categoryID := int64(queryInt(r, "category_id", 0))
	vehicles, err := h.VehicleService.GetActive(r.Context(), categoryID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Falha ao carregar veículos.", "INTERNAL_ERROR")
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"items": vehicles,
		"total": len(vehicles),
	})
}

func (h *Handler) VehicleDetails(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "identificador inválido", "INVALID_INPUT")
		return
	}
	vehicle, err := h.VehicleService.GetByID(r.Context(), id)
	if errors.Is(err, repository.ErrVehicleNotFound) {
		response.Error(w, http.StatusNotFound, "Veículo não encontrado.", "NOT_FOUND")
		return
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Falha ao carregar veículo.", "INTERNAL_ERROR")
		return
	}
	response.JSON(w, http.StatusOK, vehicle)
}

func queryInt(r *http.Request, name string, fallback int) int {
	value, err := strconv.Atoi(r.URL.Query().Get(name))
	if err != nil || value < 0 {
		return fallback
	}
	return value
}
