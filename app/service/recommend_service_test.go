package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"carro-ideal/app/models"
	"carro-ideal/app/repository"
)

type fakeQuestionRepository struct {
	questions []models.Question
}

func (r *fakeQuestionRepository) GetActive(_ context.Context) ([]models.Question, error) {
	return r.questions, nil
}

func TestScoreVehicle(t *testing.T) {
	tests := []struct {
		name           string
		userProfile    map[string]float64
		vehicleProfile map[string]float64
		want           float64
	}{
		{
			name:           "perfect match",
			userProfile:    map[string]float64{"urban": 1, "efficiency": 2},
			vehicleProfile: map[string]float64{"urban": 1, "efficiency": 1},
			want:           100,
		},
		{
			name:           "partial match",
			userProfile:    map[string]float64{"urban": 1, "efficiency": 1},
			vehicleProfile: map[string]float64{"urban": 1},
			want:           50,
		},
		{
			name:           "no match",
			userProfile:    map[string]float64{"urban": 1},
			vehicleProfile: map[string]float64{"offroad": 1},
			want:           0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			score, _ := scoreVehicle(test.userProfile, test.vehicleProfile)
			if score != test.want {
				t.Fatalf("scoreVehicle() = %.2f, want %.2f", score, test.want)
			}
		})
	}
}

func TestQuestionnaireBuildProfile(t *testing.T) {
	repo := &fakeQuestionRepository{
		questions: []models.Question{
			{
				ID:     1,
				Weight: 2,
				Options: []models.AnswerOption{
					{ID: 10, QuestionID: 1, ScoreProfile: map[string]float64{"urban": 1}},
				},
			},
			{
				ID:     2,
				Weight: 1,
				Options: []models.AnswerOption{
					{ID: 20, QuestionID: 2, ScoreProfile: map[string]float64{"efficiency": 0.5}},
				},
			},
		},
	}
	service := NewQuestionnaireService(repo)

	profile, err := service.BuildProfile(context.Background(), []models.SubmittedAnswer{
		{QuestionID: 1, AnswerOptionID: 10},
		{QuestionID: 2, AnswerOptionID: 20},
	})
	if err != nil {
		t.Fatalf("BuildProfile() error = %v", err)
	}
	if profile["urban"] != 2 || profile["efficiency"] != 0.5 {
		t.Fatalf("BuildProfile() = %#v, want weighted profile", profile)
	}
}

func TestQuestionnaireRejectsIncompleteAnswers(t *testing.T) {
	repo := &fakeQuestionRepository{
		questions: []models.Question{
			{ID: 1, Options: []models.AnswerOption{{ID: 10}}},
			{ID: 2, Options: []models.AnswerOption{{ID: 20}}},
		},
	}
	service := NewQuestionnaireService(repo)

	_, err := service.BuildProfile(context.Background(), []models.SubmittedAnswer{
		{QuestionID: 1, AnswerOptionID: 10},
	})
	if !errors.Is(err, ErrIncompleteQuestionnaire) {
		t.Fatalf("BuildProfile() error = %v, want ErrIncompleteQuestionnaire", err)
	}
}

type fakeVehicleRepository struct {
	vehicles []models.Vehicle
	err      error
}

func (r *fakeVehicleRepository) GetActive(_ context.Context, _ int64) ([]models.Vehicle, error) {
	return r.vehicles, r.err
}

func (r *fakeVehicleRepository) GetByID(_ context.Context, id int64) (*models.Vehicle, error) {
	for _, v := range r.vehicles {
		if v.ID == id {
			return &v, nil
		}
	}
	return nil, repository.ErrVehicleNotFound
}

type fakeRecommendationRepository struct {
	stored       *models.Recommendation
	history      []models.Recommendation
	historyTotal int
	byID         *models.Recommendation
	byIDErr      error
	createErr    error
}

func (r *fakeRecommendationRepository) Create(_ context.Context, rec *models.Recommendation, _ []models.SubmittedAnswer) error {
	if r.createErr != nil {
		return r.createErr
	}
	rec.ID = 1
	r.stored = rec
	return nil
}

func (r *fakeRecommendationRepository) GetByUser(_ context.Context, _ int64, limit, offset int) ([]models.Recommendation, int, error) {
	return r.history, r.historyTotal, nil
}

func (r *fakeRecommendationRepository) GetByID(_ context.Context, id, userID int64) (*models.Recommendation, error) {
	if r.byIDErr != nil {
		return nil, r.byIDErr
	}
	return r.byID, nil
}

func makeRecommendationService(questions []models.Question, vehicles []models.Vehicle, recRepo repository.RecommendationRepository) *RecommendationService {
	questionRepo := &fakeQuestionRepository{questions: questions}
	vehicleRepo := &fakeVehicleRepository{vehicles: vehicles}
	qService := NewQuestionnaireService(questionRepo)
	vService := NewVehicleService(vehicleRepo)
	return NewRecommendationService(qService, vService, recRepo)
}

func TestRecommendationServiceGenerate(t *testing.T) {
	questions := []models.Question{
		{
			ID: 1, Weight: 1,
			Options: []models.AnswerOption{
				{ID: 10, QuestionID: 1, ScoreProfile: map[string]float64{"urban": 1}},
			},
		},
	}
	vehicles := []models.Vehicle{
		{ID: 1, MatchProfile: map[string]float64{"urban": 1}},
		{ID: 2, MatchProfile: map[string]float64{"offroad": 1}},
	}
	recRepo := &fakeRecommendationRepository{}
	service := makeRecommendationService(questions, vehicles, recRepo)

	answers := []models.SubmittedAnswer{{QuestionID: 1, AnswerOptionID: 10}}
	rec, err := service.Generate(context.Background(), 42, answers)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if rec.UserID != 42 {
		t.Fatalf("Generate() userID = %d, want 42", rec.UserID)
	}
	if len(rec.Items) == 0 {
		t.Fatal("Generate() returned no items")
	}
	if rec.Items[0].Vehicle.ID != 1 {
		t.Fatalf("Generate() first item vehicle ID = %d, want 1 (highest score)", rec.Items[0].Vehicle.ID)
	}
	if rec.Items[0].Rank != 1 {
		t.Fatalf("Generate() first item rank = %d, want 1", rec.Items[0].Rank)
	}
}

func TestRecommendationServiceGenerateNoVehicles(t *testing.T) {
	questions := []models.Question{
		{ID: 1, Weight: 1, Options: []models.AnswerOption{{ID: 10, QuestionID: 1, ScoreProfile: map[string]float64{"urban": 1}}}},
	}
	service := makeRecommendationService(questions, nil, &fakeRecommendationRepository{})
	_, err := service.Generate(context.Background(), 1, []models.SubmittedAnswer{{QuestionID: 1, AnswerOptionID: 10}})
	if !errors.Is(err, ErrNoVehicles) {
		t.Fatalf("Generate() error = %v, want ErrNoVehicles", err)
	}
}

func TestRecommendationServiceHistory(t *testing.T) {
	recs := []models.Recommendation{{ID: 1, UserID: 42}, {ID: 2, UserID: 42}}
	recRepo := &fakeRecommendationRepository{history: recs, historyTotal: 2}
	service := makeRecommendationService(nil, nil, recRepo)

	items, total, err := service.History(context.Background(), 42, 1, 10)
	if err != nil {
		t.Fatalf("History() error = %v", err)
	}
	if total != 2 || len(items) != 2 {
		t.Fatalf("History() returned %d items (total %d), want 2 items (total 2)", len(items), total)
	}
}

func TestRecommendationServiceHistoryPaginationDefaults(t *testing.T) {
	recRepo := &fakeRecommendationRepository{}
	service := makeRecommendationService(nil, nil, recRepo)
	_, _, err := service.History(context.Background(), 1, 0, 0)
	if err != nil {
		t.Fatalf("History() error = %v", err)
	}
	_, _, err = service.History(context.Background(), 1, 1, 200)
	if err != nil {
		t.Fatalf("History() with large limit error = %v", err)
	}
}

func TestRecommendationServiceGet(t *testing.T) {
	rec := &models.Recommendation{ID: 5, UserID: 42}
	recRepo := &fakeRecommendationRepository{byID: rec}
	service := makeRecommendationService(nil, nil, recRepo)

	got, err := service.Get(context.Background(), 5, 42)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.ID != 5 {
		t.Fatalf("Get() ID = %d, want 5", got.ID)
	}
}

func TestRecommendationServiceGetNotFound(t *testing.T) {
	recRepo := &fakeRecommendationRepository{byIDErr: repository.ErrRecommendationNotFound}
	service := makeRecommendationService(nil, nil, recRepo)

	_, err := service.Get(context.Background(), 99, 1)
	if !errors.Is(err, repository.ErrRecommendationNotFound) {
		t.Fatalf("Get() error = %v, want ErrRecommendationNotFound", err)
	}
}

func TestBuildReason(t *testing.T) {
	tests := []struct {
		name    string
		matches []string
		wantNot string
	}{
		{name: "no matches", matches: nil, wantNot: ""},
		{name: "known dimension", matches: []string{"urban"}, wantNot: ""},
		{name: "multiple dimensions", matches: []string{"urban", "efficiency", "comfort"}, wantNot: ""},
		{name: "unknown dimension", matches: []string{"unknown_key"}, wantNot: ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reason := buildReason(tc.matches)
			if reason == "" {
				t.Fatal("buildReason() returned empty string")
			}
		})
	}
}

func TestQuestionnaireServiceGetActiveWithCache(t *testing.T) {
	questions := []models.Question{{ID: 1, Text: "Pergunta 1"}}
	repo := &fakeQuestionRepository{questions: questions}
	cache := NewCatalogCache(time.Minute)
	service := NewQuestionnaireService(repo, cache)

	got, err := service.GetActive(context.Background())
	if err != nil || len(got) != 1 {
		t.Fatalf("GetActive() error = %v, want 1 question", err)
	}
	if _, found := cache.Questions(); !found {
		t.Fatal("GetActive() should populate cache")
	}

	repo.questions = nil
	got, err = service.GetActive(context.Background())
	if err != nil || len(got) != 1 {
		t.Fatal("GetActive() should serve from cache on second call")
	}
}

func TestVehicleServiceGetActive(t *testing.T) {
	vehicles := []models.Vehicle{{ID: 1}, {ID: 2}}
	repo := &fakeVehicleRepository{vehicles: vehicles}
	cache := NewCatalogCache(time.Minute)
	service := NewVehicleService(repo, cache)

	got, err := service.GetActive(context.Background(), 0)
	if err != nil || len(got) != 2 {
		t.Fatalf("GetActive() error = %v, want 2 vehicles", err)
	}
	if _, found := cache.Vehicles(0); !found {
		t.Fatal("GetActive() should populate cache")
	}

	repo.vehicles = nil
	got, err = service.GetActive(context.Background(), 0)
	if err != nil || len(got) != 2 {
		t.Fatal("GetActive() should serve from cache on second call")
	}
}

func TestVehicleServiceGetByID(t *testing.T) {
	vehicles := []models.Vehicle{{ID: 7, Brand: "Toyota"}}
	service := NewVehicleService(&fakeVehicleRepository{vehicles: vehicles})

	v, err := service.GetByID(context.Background(), 7)
	if err != nil || v.Brand != "Toyota" {
		t.Fatalf("GetByID() error = %v, want Toyota", err)
	}

	_, err = service.GetByID(context.Background(), 999)
	if !errors.Is(err, repository.ErrVehicleNotFound) {
		t.Fatalf("GetByID() error = %v, want ErrVehicleNotFound", err)
	}
}

func TestRecommendationServiceLimitsTop10(t *testing.T) {
	questions := []models.Question{
		{ID: 1, Weight: 1, Options: []models.AnswerOption{{ID: 10, QuestionID: 1, ScoreProfile: map[string]float64{"urban": 1}}}},
	}
	vehicles := make([]models.Vehicle, 15)
	for i := range vehicles {
		vehicles[i] = models.Vehicle{ID: int64(i + 1), MatchProfile: map[string]float64{"urban": 1}}
	}
	recRepo := &fakeRecommendationRepository{}
	service := makeRecommendationService(questions, vehicles, recRepo)

	rec, err := service.Generate(context.Background(), 1, []models.SubmittedAnswer{{QuestionID: 1, AnswerOptionID: 10}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if len(rec.Items) > 10 {
		t.Fatalf("Generate() returned %d items, want at most 10", len(rec.Items))
	}
}

func TestRecommendationServiceGenerateStoreError(t *testing.T) {
	questions := []models.Question{
		{ID: 1, Weight: 1, Options: []models.AnswerOption{{ID: 10, QuestionID: 1, ScoreProfile: map[string]float64{"urban": 1}}}},
	}
	vehicles := []models.Vehicle{{ID: 1, MatchProfile: map[string]float64{"urban": 1}}}
	recRepo := &fakeRecommendationRepository{createErr: errors.New("db error")}
	service := makeRecommendationService(questions, vehicles, recRepo)

	_, err := service.Generate(context.Background(), 1, []models.SubmittedAnswer{{QuestionID: 1, AnswerOptionID: 10}})
	if err == nil {
		t.Fatal("Generate() should propagate repository errors")
	}
}

var _ = time.Now

func BenchmarkScoreVehicle(b *testing.B) {
	userProfile := map[string]float64{
		"urban": 1, "efficiency": 1, "comfort": 1, "space": 1,
		"family": 1, "performance": 1, "automatic": 1, "cargo": 1,
	}
	vehicleProfile := map[string]float64{
		"urban": 0.9, "efficiency": 0.8, "comfort": 0.7, "space": 0.6,
		"family": 0.8, "performance": 0.5, "automatic": 1, "cargo": 0.6,
	}
	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		scoreVehicle(userProfile, vehicleProfile)
	}
}
