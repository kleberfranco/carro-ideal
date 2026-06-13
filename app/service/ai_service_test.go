package service

import (
	"context"
	"errors"
	"testing"

	"carro-ideal/app/models"
)

type fakeAIClient struct {
	response string
	err      error
}

func (f *fakeAIClient) ChatComplete(_ context.Context, _, _ string) (string, error) {
	return f.response, f.err
}

var validAIResponse = `{
  "summary": "Perfil voltado para uso urbano e economia.",
  "items": [
    {"vehicle_id": 1, "rank": 1, "reason": "Excelente para cidade, câmbio automático."},
    {"vehicle_id": 2, "rank": 2, "reason": "Bom custo-benefício, econômico."}
  ]
}`

func makeAITestData() ([]models.SubmittedAnswer, []models.Question, []models.Vehicle) {
	questions := []models.Question{
		{
			ID: 1, Text: "Qual o principal uso do carro?",
			Options: []models.AnswerOption{
				{ID: 10, QuestionID: 1, Text: "Uso urbano diário"},
			},
		},
	}
	answers := []models.SubmittedAnswer{{QuestionID: 1, AnswerOptionID: 10}}
	vehicles := []models.Vehicle{
		{ID: 1, Brand: "Toyota", Model: "Corolla", Year: 2024, Category: models.VehicleCategory{Name: "Sedã"},
			PriceMin: 120000, PriceMax: 140000, Transmission: "automatic", Description: "Confortável e econômico."},
		{ID: 2, Brand: "Volkswagen", Model: "Gol", Year: 2023, Category: models.VehicleCategory{Name: "Hatchback"},
			PriceMin: 70000, PriceMax: 80000, Transmission: "manual"},
	}
	return answers, questions, vehicles
}

func TestAIServiceRecommend(t *testing.T) {
	answers, questions, vehicles := makeAITestData()
	svc := NewAIService(&fakeAIClient{response: validAIResponse})

	rec, err := svc.Recommend(context.Background(), answers, questions, vehicles)
	if err != nil {
		t.Fatalf("Recommend() error = %v", err)
	}
	if rec.Summary == "" {
		t.Fatal("Recommend() summary is empty")
	}
	if len(rec.Items) != 2 {
		t.Fatalf("Recommend() items = %d, want 2", len(rec.Items))
	}
	if rec.Items[0].VehicleID != 1 || rec.Items[0].Rank != 1 {
		t.Fatalf("Recommend() first item = vehicle %d rank %d, want vehicle 1 rank 1",
			rec.Items[0].VehicleID, rec.Items[0].Rank)
	}
}

func TestAIServiceRecommendClientError(t *testing.T) {
	answers, questions, vehicles := makeAITestData()
	svc := NewAIService(&fakeAIClient{err: errors.New("network error")})

	_, err := svc.Recommend(context.Background(), answers, questions, vehicles)
	if err == nil {
		t.Fatal("Recommend() should propagate client error")
	}
}

func TestAIServiceRecommendInvalidJSON(t *testing.T) {
	answers, questions, vehicles := makeAITestData()
	svc := NewAIService(&fakeAIClient{response: "not json"})

	_, err := svc.Recommend(context.Background(), answers, questions, vehicles)
	if err == nil {
		t.Fatal("Recommend() should error on invalid JSON")
	}
}

func TestAIServiceRecommendMarkdownFence(t *testing.T) {
	answers, questions, vehicles := makeAITestData()
	wrapped := "```json\n" + validAIResponse + "\n```"
	svc := NewAIService(&fakeAIClient{response: wrapped})

	rec, err := svc.Recommend(context.Background(), answers, questions, vehicles)
	if err != nil {
		t.Fatalf("Recommend() should strip markdown fences, got error: %v", err)
	}
	if len(rec.Items) != 2 {
		t.Fatalf("Recommend() items = %d, want 2", len(rec.Items))
	}
}

func TestAIServiceRecommendEmptyItems(t *testing.T) {
	answers, questions, vehicles := makeAITestData()
	svc := NewAIService(&fakeAIClient{response: `{"summary":"ok","items":[]}`})

	_, err := svc.Recommend(context.Background(), answers, questions, vehicles)
	if err == nil {
		t.Fatal("Recommend() should error when items is empty")
	}
}

func TestBuildPromptContainsVehicleIDs(t *testing.T) {
	_, questions, vehicles := makeAITestData()
	answers := []models.SubmittedAnswer{{QuestionID: 1, AnswerOptionID: 10}}

	prompt := buildPrompt(answers, questions, vehicles)
	if prompt == "" {
		t.Fatal("buildPrompt() returned empty string")
	}
	for _, v := range vehicles {
		expected := "Toyota Corolla"
		if v.Brand == "Volkswagen" {
			expected = "Volkswagen Gol"
		}
		_ = expected
	}
	// Verify vehicle IDs are present for reliable matching
	if !contains(prompt, "[ID:1]") || !contains(prompt, "[ID:2]") {
		t.Fatal("buildPrompt() should include vehicle IDs in the prompt")
	}
	// Verify question text is present
	if !contains(prompt, "Qual o principal uso do carro?") {
		t.Fatal("buildPrompt() should include question text")
	}
	// Verify answer text is present
	if !contains(prompt, "Uso urbano diário") {
		t.Fatal("buildPrompt() should include selected answer text")
	}
	// Verify description is present
	if !contains(prompt, "Confortável e econômico.") {
		t.Fatal("buildPrompt() should include vehicle description")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
