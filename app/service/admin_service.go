package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"carro-ideal/app/models"
	"carro-ideal/app/repository"
)

var ErrAdminValidation = errors.New("dados administrativos inválidos")

type AdminService struct {
	repo  repository.AdminRepositoryInterface
	cache *CatalogCache
}

func NewAdminService(repo repository.AdminRepositoryInterface, caches ...*CatalogCache) *AdminService {
	var cache *CatalogCache
	if len(caches) > 0 {
		cache = caches[0]
	}
	return &AdminService{repo: repo, cache: cache}
}

func (s *AdminService) Stats(ctx context.Context) (*models.AdminStats, error) {
	return s.repo.Stats(ctx)
}

func (s *AdminService) Users(ctx context.Context, search string, page, limit int) ([]models.User, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.Users(ctx, strings.TrimSpace(search), limit, (page-1)*limit)
}

func (s *AdminService) Vehicles(ctx context.Context, search string, page, limit int) ([]models.Vehicle, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.Vehicles(ctx, strings.TrimSpace(search), limit, (page-1)*limit)
}

func (s *AdminService) CreateVehicle(ctx context.Context, input models.VehicleInput) (*models.Vehicle, error) {
	vehicle := &models.Vehicle{Active: true, MatchProfile: input.MatchProfile}
	applyVehicleInput(vehicle, input)
	if err := validateVehicle(vehicle); err != nil {
		return nil, err
	}
	if err := s.repo.CreateVehicle(ctx, vehicle); err != nil {
		return nil, err
	}
	s.invalidateCatalog()
	return s.repo.Vehicle(ctx, vehicle.ID)
}

func (s *AdminService) UpdateVehicle(ctx context.Context, id int64, input models.VehicleInput) (*models.Vehicle, error) {
	vehicle, err := s.repo.Vehicle(ctx, id)
	if err != nil {
		return nil, err
	}
	applyVehicleInput(vehicle, input)
	if err := validateVehicle(vehicle); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateVehicle(ctx, vehicle); err != nil {
		return nil, err
	}
	s.invalidateCatalog()
	return s.repo.Vehicle(ctx, id)
}

func (s *AdminService) DeleteVehicle(ctx context.Context, id int64) error {
	err := s.repo.DeleteVehicle(ctx, id)
	if err == nil {
		s.invalidateCatalog()
	}
	return err
}

func (s *AdminService) Categories(ctx context.Context) ([]models.VehicleCategory, error) {
	return s.repo.Categories(ctx)
}

func (s *AdminService) CreateCategory(ctx context.Context, input models.CategoryInput) (*models.VehicleCategory, error) {
	category := &models.VehicleCategory{Active: true}
	applyCategoryInput(category, input)
	if strings.TrimSpace(category.Name) == "" {
		return nil, ErrAdminValidation
	}
	if err := s.repo.CreateCategory(ctx, category); err != nil {
		return nil, err
	}
	s.invalidateCatalog()
	return s.repo.Category(ctx, category.ID)
}

func (s *AdminService) UpdateCategory(ctx context.Context, id int64, input models.CategoryInput) (*models.VehicleCategory, error) {
	category, err := s.repo.Category(ctx, id)
	if err != nil {
		return nil, err
	}
	applyCategoryInput(category, input)
	if strings.TrimSpace(category.Name) == "" {
		return nil, ErrAdminValidation
	}
	if err := s.repo.UpdateCategory(ctx, category); err != nil {
		return nil, err
	}
	s.invalidateCatalog()
	return s.repo.Category(ctx, id)
}

func (s *AdminService) DeleteCategory(ctx context.Context, id int64) error {
	err := s.repo.DeleteCategory(ctx, id)
	if err == nil {
		s.invalidateCatalog()
	}
	return err
}

func (s *AdminService) Questions(ctx context.Context) ([]models.Question, error) {
	return s.repo.Questions(ctx)
}

func (s *AdminService) CreateQuestion(ctx context.Context, input models.QuestionInput) (*models.Question, error) {
	question := &models.Question{Type: "SINGLE_CHOICE", Weight: 1, Active: true}
	applyQuestionInput(question, input)
	if err := validateQuestion(question); err != nil {
		return nil, err
	}
	if err := s.repo.CreateQuestion(ctx, question); err != nil {
		return nil, err
	}
	s.invalidateCatalog()
	return s.repo.Question(ctx, question.ID)
}

func (s *AdminService) UpdateQuestion(ctx context.Context, id int64, input models.QuestionInput) (*models.Question, error) {
	question, err := s.repo.Question(ctx, id)
	if err != nil {
		return nil, err
	}
	applyQuestionInput(question, input)
	if err := validateQuestion(question); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateQuestion(ctx, question); err != nil {
		return nil, err
	}
	s.invalidateCatalog()
	return s.repo.Question(ctx, id)
}

func (s *AdminService) DeleteQuestion(ctx context.Context, id int64) error {
	err := s.repo.DeleteQuestion(ctx, id)
	if err == nil {
		s.invalidateCatalog()
	}
	return err
}

func (s *AdminService) CreateOption(ctx context.Context, questionID int64, input models.AnswerOptionInput) (*models.AnswerOption, error) {
	if _, err := s.repo.Question(ctx, questionID); err != nil {
		return nil, err
	}
	option := &models.AnswerOption{
		QuestionID:   questionID,
		ScoreProfile: input.ScoreProfile,
		Active:       true,
	}
	applyOptionInput(option, input)
	if err := validateOption(option); err != nil {
		return nil, err
	}
	if err := s.repo.CreateOption(ctx, option); err != nil {
		return nil, err
	}
	s.invalidateCatalog()
	return option, nil
}

func (s *AdminService) UpdateOption(ctx context.Context, questionID, optionID int64, input models.AnswerOptionInput) (*models.AnswerOption, error) {
	question, err := s.repo.Question(ctx, questionID)
	if err != nil {
		return nil, err
	}
	var option *models.AnswerOption
	for index := range question.Options {
		if question.Options[index].ID == optionID {
			copy := question.Options[index]
			option = &copy
			break
		}
	}
	if option == nil {
		return nil, repository.ErrOptionNotFound
	}
	applyOptionInput(option, input)
	if err := validateOption(option); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateOption(ctx, option); err != nil {
		return nil, err
	}
	s.invalidateCatalog()
	return option, nil
}

func (s *AdminService) DeleteOption(ctx context.Context, questionID, optionID int64) error {
	err := s.repo.DeleteOption(ctx, questionID, optionID)
	if err == nil {
		s.invalidateCatalog()
	}
	return err
}

func (s *AdminService) invalidateCatalog() {
	if s.cache != nil {
		s.cache.Invalidate()
	}
}

func applyVehicleInput(vehicle *models.Vehicle, input models.VehicleInput) {
	if input.CategoryID != nil {
		vehicle.CategoryID = *input.CategoryID
	}
	if input.Brand != nil {
		vehicle.Brand = strings.TrimSpace(*input.Brand)
	}
	if input.Model != nil {
		vehicle.Model = strings.TrimSpace(*input.Model)
	}
	if input.Version != nil {
		vehicle.Version = strings.TrimSpace(*input.Version)
	}
	if input.Year != nil {
		vehicle.Year = *input.Year
	}
	if input.FuelType != nil {
		vehicle.FuelType = strings.TrimSpace(*input.FuelType)
	}
	if input.Transmission != nil {
		vehicle.Transmission = strings.TrimSpace(*input.Transmission)
	}
	if input.PriceMin != nil {
		vehicle.PriceMin = *input.PriceMin
	}
	if input.PriceMax != nil {
		vehicle.PriceMax = *input.PriceMax
	}
	if input.Seats != nil {
		vehicle.Seats = *input.Seats
	}
	if input.TrunkCapacity != nil {
		vehicle.TrunkCapacity = *input.TrunkCapacity
	}
	if input.ConsumptionCity != nil {
		vehicle.ConsumptionCity = *input.ConsumptionCity
	}
	if input.ConsumptionHighway != nil {
		vehicle.ConsumptionHighway = *input.ConsumptionHighway
	}
	if input.Description != nil {
		vehicle.Description = strings.TrimSpace(*input.Description)
	}
	if input.Strengths != nil {
		vehicle.Strengths = strings.TrimSpace(*input.Strengths)
	}
	if input.Weaknesses != nil {
		vehicle.Weaknesses = strings.TrimSpace(*input.Weaknesses)
	}
	if input.MatchProfile != nil {
		vehicle.MatchProfile = input.MatchProfile
	}
	if input.Active != nil {
		vehicle.Active = *input.Active
	}
}

func validateVehicle(vehicle *models.Vehicle) error {
	currentYear := time.Now().Year() + 1
	if vehicle.CategoryID < 1 || vehicle.Brand == "" || vehicle.Model == "" ||
		vehicle.FuelType == "" || vehicle.Transmission == "" ||
		vehicle.Year < 1980 || vehicle.Year > currentYear ||
		vehicle.PriceMin < 0 || vehicle.PriceMax < vehicle.PriceMin ||
		len(vehicle.MatchProfile) == 0 {
		return ErrAdminValidation
	}
	return nil
}

func applyCategoryInput(category *models.VehicleCategory, input models.CategoryInput) {
	if input.Name != nil {
		category.Name = strings.TrimSpace(*input.Name)
	}
	if input.Description != nil {
		category.Description = strings.TrimSpace(*input.Description)
	}
	if input.Active != nil {
		category.Active = *input.Active
	}
}

func applyQuestionInput(question *models.Question, input models.QuestionInput) {
	if input.Text != nil {
		question.Text = strings.TrimSpace(*input.Text)
	}
	if input.Type != nil {
		question.Type = strings.TrimSpace(*input.Type)
	}
	if input.Weight != nil {
		question.Weight = *input.Weight
	}
	if input.DisplayOrder != nil {
		question.DisplayOrder = *input.DisplayOrder
	}
	if input.Active != nil {
		question.Active = *input.Active
	}
}

func validateQuestion(question *models.Question) error {
	if question.Text == "" || question.Type == "" || question.Weight <= 0 || question.DisplayOrder < 1 {
		return ErrAdminValidation
	}
	return nil
}

func applyOptionInput(option *models.AnswerOption, input models.AnswerOptionInput) {
	if input.Text != nil {
		option.Text = strings.TrimSpace(*input.Text)
	}
	if input.ScoreProfile != nil {
		option.ScoreProfile = input.ScoreProfile
	}
	if input.DisplayOrder != nil {
		option.DisplayOrder = *input.DisplayOrder
	}
	if input.Active != nil {
		option.Active = *input.Active
	}
}

func validateOption(option *models.AnswerOption) error {
	if option.Text == "" || option.DisplayOrder < 1 || len(option.ScoreProfile) == 0 {
		return ErrAdminValidation
	}
	for _, value := range option.ScoreProfile {
		if value < 0 || value > 1 {
			return ErrAdminValidation
		}
	}
	return nil
}
