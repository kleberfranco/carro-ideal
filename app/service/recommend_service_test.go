package service

import (
	"context"
	"errors"
	"testing"

	"carro-ideal/app/models"
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
