package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"carro-ideal/app/models"
	"carro-ideal/app/repository"
)

type fakeAdminRepository struct {
	users      []models.User
	vehicles   []models.Vehicle
	vehicle    *models.Vehicle
	categories []models.VehicleCategory
	category   *models.VehicleCategory
	questions  []models.Question
	question   *models.Question

	createVehicleErr  error
	vehicleErr        error
	categoryErr       error
	questionErr       error
	optionNotFoundErr bool
}

func (r *fakeAdminRepository) Stats(_ context.Context) (*models.AdminStats, error) {
	return &models.AdminStats{Users: 10, Vehicles: 5}, nil
}

func (r *fakeAdminRepository) Users(_ context.Context, _ string, limit, offset int) ([]models.User, int, error) {
	return r.users, len(r.users), nil
}

func (r *fakeAdminRepository) Vehicles(_ context.Context, _ string, limit, offset int) ([]models.Vehicle, int, error) {
	return r.vehicles, len(r.vehicles), nil
}

func (r *fakeAdminRepository) Vehicle(_ context.Context, id int64) (*models.Vehicle, error) {
	if r.vehicleErr != nil {
		return nil, r.vehicleErr
	}
	if r.vehicle != nil {
		return r.vehicle, nil
	}
	for _, v := range r.vehicles {
		if v.ID == id {
			cp := v
			return &cp, nil
		}
	}
	return nil, repository.ErrVehicleNotFound
}

func (r *fakeAdminRepository) CreateVehicle(_ context.Context, vehicle *models.Vehicle) error {
	if r.createVehicleErr != nil {
		return r.createVehicleErr
	}
	vehicle.ID = 1
	r.vehicle = vehicle
	return nil
}

func (r *fakeAdminRepository) UpdateVehicle(_ context.Context, vehicle *models.Vehicle) error {
	r.vehicle = vehicle
	return nil
}

func (r *fakeAdminRepository) DeleteVehicle(_ context.Context, _ int64) error { return nil }

func (r *fakeAdminRepository) Categories(_ context.Context) ([]models.VehicleCategory, error) {
	return r.categories, nil
}

func (r *fakeAdminRepository) Category(_ context.Context, id int64) (*models.VehicleCategory, error) {
	if r.categoryErr != nil {
		return nil, r.categoryErr
	}
	if r.category != nil {
		return r.category, nil
	}
	return nil, repository.ErrCategoryNotFound
}

func (r *fakeAdminRepository) CreateCategory(_ context.Context, category *models.VehicleCategory) error {
	category.ID = 1
	r.category = category
	return nil
}

func (r *fakeAdminRepository) UpdateCategory(_ context.Context, category *models.VehicleCategory) error {
	r.category = category
	return nil
}

func (r *fakeAdminRepository) DeleteCategory(_ context.Context, _ int64) error { return nil }

func (r *fakeAdminRepository) Questions(_ context.Context) ([]models.Question, error) {
	return r.questions, nil
}

func (r *fakeAdminRepository) Question(_ context.Context, id int64) (*models.Question, error) {
	if r.questionErr != nil {
		return nil, r.questionErr
	}
	if r.question != nil {
		return r.question, nil
	}
	return nil, repository.ErrQuestionNotFound
}

func (r *fakeAdminRepository) CreateQuestion(_ context.Context, question *models.Question) error {
	question.ID = 1
	r.question = question
	return nil
}

func (r *fakeAdminRepository) UpdateQuestion(_ context.Context, question *models.Question) error {
	r.question = question
	return nil
}

func (r *fakeAdminRepository) DeleteQuestion(_ context.Context, _ int64) error { return nil }

func (r *fakeAdminRepository) CreateOption(_ context.Context, option *models.AnswerOption) error {
	option.ID = 1
	return nil
}

func (r *fakeAdminRepository) UpdateOption(_ context.Context, _ *models.AnswerOption) error {
	return nil
}

func (r *fakeAdminRepository) DeleteOption(_ context.Context, _, _ int64) error { return nil }

func TestValidateVehicle(t *testing.T) {
	valid := &models.Vehicle{
		CategoryID:   1,
		Brand:        "Toyota",
		Model:        "Corolla",
		Year:         2025,
		FuelType:     "Flex",
		Transmission: "CVT",
		PriceMin:     100000,
		PriceMax:     150000,
		MatchProfile: map[string]float64{"comfort": 1},
	}
	if err := validateVehicle(valid); err != nil {
		t.Fatalf("validateVehicle() error = %v", err)
	}

	invalid := *valid
	invalid.PriceMax = 90000
	if err := validateVehicle(&invalid); !errors.Is(err, ErrAdminValidation) {
		t.Fatalf("validateVehicle() error = %v, want ErrAdminValidation", err)
	}
}

func TestValidateOption(t *testing.T) {
	option := &models.AnswerOption{
		Text:         "Muito importante",
		DisplayOrder: 1,
		ScoreProfile: map[string]float64{"comfort": 1},
	}
	if err := validateOption(option); err != nil {
		t.Fatalf("validateOption() error = %v", err)
	}

	option.ScoreProfile["comfort"] = 1.5
	if err := validateOption(option); !errors.Is(err, ErrAdminValidation) {
		t.Fatalf("validateOption() error = %v, want ErrAdminValidation", err)
	}
}

func TestValidateVehicleRequiredFields(t *testing.T) {
	year := time.Now().Year()
	cases := []struct {
		name    string
		vehicle models.Vehicle
	}{
		{name: "missing brand", vehicle: models.Vehicle{CategoryID: 1, Model: "X", Year: year, FuelType: "Flex", Transmission: "CVT", PriceMax: 1, MatchProfile: map[string]float64{"a": 1}}},
		{name: "missing model", vehicle: models.Vehicle{CategoryID: 1, Brand: "B", Year: year, FuelType: "Flex", Transmission: "CVT", PriceMax: 1, MatchProfile: map[string]float64{"a": 1}}},
		{name: "year too old", vehicle: models.Vehicle{CategoryID: 1, Brand: "B", Model: "X", Year: 1900, FuelType: "Flex", Transmission: "CVT", PriceMax: 1, MatchProfile: map[string]float64{"a": 1}}},
		{name: "empty match profile", vehicle: models.Vehicle{CategoryID: 1, Brand: "B", Model: "X", Year: year, FuelType: "Flex", Transmission: "CVT", PriceMax: 1}},
		{name: "no category", vehicle: models.Vehicle{Brand: "B", Model: "X", Year: year, FuelType: "Flex", Transmission: "CVT", PriceMax: 1, MatchProfile: map[string]float64{"a": 1}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := validateVehicle(&tc.vehicle); !errors.Is(err, ErrAdminValidation) {
				t.Fatalf("validateVehicle() error = %v, want ErrAdminValidation", err)
			}
		})
	}
}

func TestValidateQuestion(t *testing.T) {
	valid := &models.Question{Text: "Qual o seu uso?", Type: "SINGLE_CHOICE", Weight: 1, DisplayOrder: 1}
	if err := validateQuestion(valid); err != nil {
		t.Fatalf("validateQuestion() error = %v", err)
	}

	cases := []struct {
		name     string
		question models.Question
	}{
		{name: "empty text", question: models.Question{Type: "SINGLE_CHOICE", Weight: 1, DisplayOrder: 1}},
		{name: "empty type", question: models.Question{Text: "Q", Weight: 1, DisplayOrder: 1}},
		{name: "zero weight", question: models.Question{Text: "Q", Type: "SINGLE_CHOICE", DisplayOrder: 1}},
		{name: "zero display order", question: models.Question{Text: "Q", Type: "SINGLE_CHOICE", Weight: 1}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := validateQuestion(&tc.question); !errors.Is(err, ErrAdminValidation) {
				t.Fatalf("validateQuestion() error = %v, want ErrAdminValidation", err)
			}
		})
	}
}

func TestApplyVehicleInput(t *testing.T) {
	vehicle := &models.Vehicle{}
	brand := "Honda"
	model := "  Civic  "
	year := 2025
	priceMin := float64(80000)
	priceMax := float64(120000)
	input := models.VehicleInput{
		Brand:    &brand,
		Model:    &model,
		Year:     &year,
		PriceMin: &priceMin,
		PriceMax: &priceMax,
	}
	applyVehicleInput(vehicle, input)

	if vehicle.Brand != "Honda" {
		t.Fatalf("applyVehicleInput() Brand = %q, want Honda", vehicle.Brand)
	}
	if vehicle.Model != "Civic" {
		t.Fatalf("applyVehicleInput() Model = %q, want Civic (trimmed)", vehicle.Model)
	}
	if vehicle.Year != 2025 {
		t.Fatalf("applyVehicleInput() Year = %d, want 2025", vehicle.Year)
	}
}

func TestApplyCategoryInput(t *testing.T) {
	category := &models.VehicleCategory{}
	name := "  SUV  "
	desc := "Sport Utility"
	active := true
	applyCategoryInput(category, models.CategoryInput{Name: &name, Description: &desc, Active: &active})
	if category.Name != "SUV" {
		t.Fatalf("applyCategoryInput() Name = %q, want SUV", category.Name)
	}
	if !category.Active {
		t.Fatal("applyCategoryInput() Active should be true")
	}
}

func TestApplyQuestionInput(t *testing.T) {
	question := &models.Question{}
	text := "  Qual o seu orçamento?  "
	qtype := "SINGLE_CHOICE"
	weight := 0.5
	order := 3
	active := false
	applyQuestionInput(question, models.QuestionInput{Text: &text, Type: &qtype, Weight: &weight, DisplayOrder: &order, Active: &active})
	if question.Text != "Qual o seu orçamento?" {
		t.Fatalf("applyQuestionInput() Text = %q, want trimmed", question.Text)
	}
	if question.Weight != 0.5 {
		t.Fatalf("applyQuestionInput() Weight = %f, want 0.5", question.Weight)
	}
	if question.DisplayOrder != 3 {
		t.Fatalf("applyQuestionInput() DisplayOrder = %d, want 3", question.DisplayOrder)
	}
}

func TestApplyOptionInput(t *testing.T) {
	option := &models.AnswerOption{}
	text := "  Muito importante  "
	profile := map[string]float64{"comfort": 0.8}
	order := 2
	active := true
	applyOptionInput(option, models.AnswerOptionInput{Text: &text, ScoreProfile: profile, DisplayOrder: &order, Active: &active})
	if !strings.Contains(option.Text, "Muito importante") {
		t.Fatalf("applyOptionInput() Text = %q, want trimmed", option.Text)
	}
	if option.ScoreProfile["comfort"] != 0.8 {
		t.Fatalf("applyOptionInput() ScoreProfile comfort = %f, want 0.8", option.ScoreProfile["comfort"])
	}
}

func validVehicleInput() models.VehicleInput {
	year := time.Now().Year()
	brand := "Toyota"
	model := "Corolla"
	fuelType := "Flex"
	transmission := "CVT"
	priceMin := float64(100000)
	priceMax := float64(150000)
	catID := int64(1)
	return models.VehicleInput{
		CategoryID:   &catID,
		Brand:        &brand,
		Model:        &model,
		Year:         &year,
		FuelType:     &fuelType,
		Transmission: &transmission,
		PriceMin:     &priceMin,
		PriceMax:     &priceMax,
		MatchProfile: map[string]float64{"comfort": 0.8},
	}
}

func TestAdminServiceStats(t *testing.T) {
	service := NewAdminService(&fakeAdminRepository{})
	stats, err := service.Stats(context.Background())
	if err != nil {
		t.Fatalf("Stats() error = %v", err)
	}
	if stats.Users != 10 {
		t.Fatalf("Stats() users = %d, want 10", stats.Users)
	}
}

func TestAdminServiceVehicles(t *testing.T) {
	repo := &fakeAdminRepository{vehicles: []models.Vehicle{{ID: 1}, {ID: 2}}}
	service := NewAdminService(repo)
	vehicles, total, err := service.Vehicles(context.Background(), "", 10, 1)
	if err != nil || total != 2 || len(vehicles) != 2 {
		t.Fatalf("Vehicles() error = %v, total = %d, len = %d", err, total, len(vehicles))
	}
}

func TestAdminServiceCreateVehicleValid(t *testing.T) {
	repo := &fakeAdminRepository{}
	service := NewAdminService(repo)
	input := validVehicleInput()

	repo.vehicle = &models.Vehicle{
		ID: 1, CategoryID: 1, Brand: "Toyota", Model: "Corolla",
		Year: time.Now().Year(), FuelType: "Flex", Transmission: "CVT",
		PriceMin: 100000, PriceMax: 150000, Active: true,
		MatchProfile: map[string]float64{"comfort": 0.8},
	}

	v, err := service.CreateVehicle(context.Background(), input)
	if err != nil {
		t.Fatalf("CreateVehicle() error = %v", err)
	}
	if v.Brand != "Toyota" {
		t.Fatalf("CreateVehicle() brand = %q, want Toyota", v.Brand)
	}
}

func TestAdminServiceCreateVehicleValidationError(t *testing.T) {
	service := NewAdminService(&fakeAdminRepository{})
	input := models.VehicleInput{}
	_, err := service.CreateVehicle(context.Background(), input)
	if !errors.Is(err, ErrAdminValidation) {
		t.Fatalf("CreateVehicle() error = %v, want ErrAdminValidation", err)
	}
}

func TestAdminServiceUpdateVehicle(t *testing.T) {
	year := time.Now().Year()
	existing := &models.Vehicle{
		ID: 1, CategoryID: 1, Brand: "Toyota", Model: "Corolla",
		Year: year, FuelType: "Flex", Transmission: "CVT",
		PriceMin: 100000, PriceMax: 150000, Active: true,
		MatchProfile: map[string]float64{"comfort": 0.8},
	}
	repo := &fakeAdminRepository{vehicle: existing}
	service := NewAdminService(repo)

	newBrand := "Honda"
	_, err := service.UpdateVehicle(context.Background(), 1, models.VehicleInput{Brand: &newBrand})
	if err != nil {
		t.Fatalf("UpdateVehicle() error = %v", err)
	}
}

func TestAdminServiceDeleteVehicle(t *testing.T) {
	service := NewAdminService(&fakeAdminRepository{})
	if err := service.DeleteVehicle(context.Background(), 1); err != nil {
		t.Fatalf("DeleteVehicle() error = %v", err)
	}
}

func TestAdminServiceCategories(t *testing.T) {
	repo := &fakeAdminRepository{categories: []models.VehicleCategory{{ID: 1, Name: "SUV"}}}
	service := NewAdminService(repo)
	cats, err := service.Categories(context.Background())
	if err != nil || len(cats) != 1 {
		t.Fatalf("Categories() error = %v, len = %d", err, len(cats))
	}
}

func TestAdminServiceCreateCategory(t *testing.T) {
	repo := &fakeAdminRepository{}
	service := NewAdminService(repo)
	name := "Sedã"
	cat, err := service.CreateCategory(context.Background(), models.CategoryInput{Name: &name})
	if err != nil {
		t.Fatalf("CreateCategory() error = %v", err)
	}
	if cat.Name != "Sedã" {
		t.Fatalf("CreateCategory() name = %q, want Sedã", cat.Name)
	}
}

func TestAdminServiceCreateCategoryEmptyName(t *testing.T) {
	service := NewAdminService(&fakeAdminRepository{})
	name := "   "
	_, err := service.CreateCategory(context.Background(), models.CategoryInput{Name: &name})
	if !errors.Is(err, ErrAdminValidation) {
		t.Fatalf("CreateCategory() error = %v, want ErrAdminValidation", err)
	}
}

func TestAdminServiceUpdateCategory(t *testing.T) {
	repo := &fakeAdminRepository{category: &models.VehicleCategory{ID: 1, Name: "SUV"}}
	service := NewAdminService(repo)
	name := "SUV Atualizado"
	_, err := service.UpdateCategory(context.Background(), 1, models.CategoryInput{Name: &name})
	if err != nil {
		t.Fatalf("UpdateCategory() error = %v", err)
	}
}

func TestAdminServiceDeleteCategory(t *testing.T) {
	service := NewAdminService(&fakeAdminRepository{})
	if err := service.DeleteCategory(context.Background(), 1); err != nil {
		t.Fatalf("DeleteCategory() error = %v", err)
	}
}

func TestAdminServiceQuestions(t *testing.T) {
	repo := &fakeAdminRepository{questions: []models.Question{{ID: 1, Text: "Qual o uso?"}}}
	service := NewAdminService(repo)
	qs, err := service.Questions(context.Background())
	if err != nil || len(qs) != 1 {
		t.Fatalf("Questions() error = %v, len = %d", err, len(qs))
	}
}

func TestAdminServiceCreateQuestion(t *testing.T) {
	repo := &fakeAdminRepository{}
	service := NewAdminService(repo)
	text := "Qual o seu orçamento?"
	qtype := "SINGLE_CHOICE"
	weight := 1.0
	order := 1
	q, err := service.CreateQuestion(context.Background(), models.QuestionInput{
		Text: &text, Type: &qtype, Weight: &weight, DisplayOrder: &order,
	})
	if err != nil {
		t.Fatalf("CreateQuestion() error = %v", err)
	}
	if q.Text != text {
		t.Fatalf("CreateQuestion() text = %q, want %q", q.Text, text)
	}
}

func TestAdminServiceCreateQuestionValidationError(t *testing.T) {
	service := NewAdminService(&fakeAdminRepository{})
	_, err := service.CreateQuestion(context.Background(), models.QuestionInput{})
	if !errors.Is(err, ErrAdminValidation) {
		t.Fatalf("CreateQuestion() error = %v, want ErrAdminValidation", err)
	}
}

func TestAdminServiceUpdateQuestion(t *testing.T) {
	existing := &models.Question{ID: 1, Text: "Pergunta", Type: "SINGLE_CHOICE", Weight: 1, DisplayOrder: 1, Active: true}
	repo := &fakeAdminRepository{question: existing}
	service := NewAdminService(repo)
	newText := "Pergunta atualizada"
	_, err := service.UpdateQuestion(context.Background(), 1, models.QuestionInput{Text: &newText})
	if err != nil {
		t.Fatalf("UpdateQuestion() error = %v", err)
	}
}

func TestAdminServiceDeleteQuestion(t *testing.T) {
	service := NewAdminService(&fakeAdminRepository{})
	if err := service.DeleteQuestion(context.Background(), 1); err != nil {
		t.Fatalf("DeleteQuestion() error = %v", err)
	}
}

func TestAdminServiceCreateOption(t *testing.T) {
	repo := &fakeAdminRepository{question: &models.Question{ID: 1, Text: "Q", Type: "SINGLE_CHOICE", Weight: 1, DisplayOrder: 1}}
	service := NewAdminService(repo)
	text := "Opção A"
	order := 1
	opt, err := service.CreateOption(context.Background(), 1, models.AnswerOptionInput{
		Text:         &text,
		ScoreProfile: map[string]float64{"comfort": 0.8},
		DisplayOrder: &order,
	})
	if err != nil {
		t.Fatalf("CreateOption() error = %v", err)
	}
	if opt.Text != "Opção A" {
		t.Fatalf("CreateOption() text = %q, want Opção A", opt.Text)
	}
}

func TestAdminServiceCreateOptionQuestionNotFound(t *testing.T) {
	repo := &fakeAdminRepository{questionErr: repository.ErrQuestionNotFound}
	service := NewAdminService(repo)
	text := "X"
	order := 1
	_, err := service.CreateOption(context.Background(), 99, models.AnswerOptionInput{
		Text: &text, ScoreProfile: map[string]float64{"a": 0.5}, DisplayOrder: &order,
	})
	if !errors.Is(err, repository.ErrQuestionNotFound) {
		t.Fatalf("CreateOption() error = %v, want ErrQuestionNotFound", err)
	}
}

func TestAdminServiceUpdateOption(t *testing.T) {
	opt := models.AnswerOption{ID: 10, QuestionID: 1, Text: "Opt", DisplayOrder: 1, ScoreProfile: map[string]float64{"a": 0.5}}
	q := &models.Question{ID: 1, Text: "Q", Options: []models.AnswerOption{opt}}
	repo := &fakeAdminRepository{question: q}
	service := NewAdminService(repo)
	newText := "Opt Atualizado"
	order := 1
	_, err := service.UpdateOption(context.Background(), 1, 10, models.AnswerOptionInput{
		Text: &newText, ScoreProfile: map[string]float64{"a": 0.5}, DisplayOrder: &order,
	})
	if err != nil {
		t.Fatalf("UpdateOption() error = %v", err)
	}
}

func TestAdminServiceUpdateOptionNotFound(t *testing.T) {
	q := &models.Question{ID: 1, Text: "Q", Options: []models.AnswerOption{}}
	repo := &fakeAdminRepository{question: q}
	service := NewAdminService(repo)
	text := "X"
	order := 1
	_, err := service.UpdateOption(context.Background(), 1, 99, models.AnswerOptionInput{
		Text: &text, ScoreProfile: map[string]float64{"a": 0.5}, DisplayOrder: &order,
	})
	if !errors.Is(err, repository.ErrOptionNotFound) {
		t.Fatalf("UpdateOption() error = %v, want ErrOptionNotFound", err)
	}
}

func TestAdminServiceDeleteOption(t *testing.T) {
	service := NewAdminService(&fakeAdminRepository{})
	if err := service.DeleteOption(context.Background(), 1, 1); err != nil {
		t.Fatalf("DeleteOption() error = %v", err)
	}
}

func TestAdminServiceInvalidateCatalogWithCache(t *testing.T) {
	cache := NewCatalogCache(time.Minute)
	cache.SetVehicles(0, []models.Vehicle{{ID: 1}})
	service := NewAdminService(&fakeAdminRepository{}, cache)
	if err := service.DeleteVehicle(context.Background(), 1); err != nil {
		t.Fatalf("DeleteVehicle() error = %v", err)
	}
	if _, found := cache.Vehicles(0); found {
		t.Fatal("cache should be invalidated after DeleteVehicle")
	}
}

var _ = time.Now
