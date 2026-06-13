package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"

	"carro-ideal/app/models"
	"carro-ideal/app/repository"
)

var (
	ErrIncompleteQuestionnaire = errors.New("responda todas as perguntas do questionário")
	ErrInvalidAnswer           = errors.New("resposta inválida para o questionário")
	ErrNoVehicles              = errors.New("nenhum veículo disponível para recomendação")
)

type QuestionnaireService struct {
	repo  repository.QuestionRepository
	cache *CatalogCache
}

func NewQuestionnaireService(repo repository.QuestionRepository, caches ...*CatalogCache) *QuestionnaireService {
	var cache *CatalogCache
	if len(caches) > 0 {
		cache = caches[0]
	}
	return &QuestionnaireService{repo: repo, cache: cache}
}

func (s *QuestionnaireService) GetActive(ctx context.Context) ([]models.Question, error) {
	if s.cache != nil {
		if questions, found := s.cache.Questions(); found {
			return questions, nil
		}
	}
	questions, err := s.repo.GetActive(ctx)
	if err == nil && s.cache != nil {
		s.cache.SetQuestions(questions)
	}
	return questions, err
}

func (s *QuestionnaireService) BuildProfile(ctx context.Context, answers []models.SubmittedAnswer) (map[string]float64, error) {
	questions, err := s.GetActive(ctx)
	if err != nil {
		return nil, err
	}
	if len(answers) != len(questions) {
		return nil, ErrIncompleteQuestionnaire
	}

	selected := make(map[int64]int64, len(answers))
	for _, answer := range answers {
		if _, duplicate := selected[answer.QuestionID]; duplicate {
			return nil, ErrInvalidAnswer
		}
		selected[answer.QuestionID] = answer.AnswerOptionID
	}

	profile := map[string]float64{}
	for _, question := range questions {
		optionID, answered := selected[question.ID]
		if !answered {
			return nil, ErrIncompleteQuestionnaire
		}

		var selectedOption *models.AnswerOption
		for index := range question.Options {
			if question.Options[index].ID == optionID {
				selectedOption = &question.Options[index]
				break
			}
		}
		if selectedOption == nil {
			return nil, ErrInvalidAnswer
		}

		for dimension, value := range selectedOption.ScoreProfile {
			profile[dimension] += value * question.Weight
		}
	}
	return profile, nil
}

type RecommendationService struct {
	questionnaire      *QuestionnaireService
	vehicles           *VehicleService
	recommendationRepo repository.RecommendationRepository
	ai                 *AIService
}

func NewRecommendationService(
	questionnaire *QuestionnaireService,
	vehicles *VehicleService,
	recommendationRepo repository.RecommendationRepository,
	aiSvc ...*AIService,
) *RecommendationService {
	var ai *AIService
	if len(aiSvc) > 0 {
		ai = aiSvc[0]
	}
	return &RecommendationService{
		questionnaire:      questionnaire,
		vehicles:           vehicles,
		recommendationRepo: recommendationRepo,
		ai:                 ai,
	}
}

func (s *RecommendationService) Generate(ctx context.Context, userID int64, answers []models.SubmittedAnswer) (*models.Recommendation, error) {
	userProfile, err := s.questionnaire.BuildProfile(ctx, answers)
	if err != nil {
		return nil, err
	}

	vehicles, err := s.vehicles.GetActive(ctx, 0)
	if err != nil {
		return nil, err
	}
	if len(vehicles) == 0 {
		return nil, ErrNoVehicles
	}

	// Try ChatGPT first; fall back to scoring algorithm on any error.
	if s.ai != nil {
		questions, qErr := s.questionnaire.GetActive(ctx)
		if qErr == nil {
			if rec, aErr := s.generateWithAI(ctx, userID, answers, questions, vehicles); aErr == nil {
				if err := s.recommendationRepo.Create(ctx, rec, answers); err != nil {
					return nil, err
				}
				return rec, nil
			}
		}
	}

	return s.generateWithScoring(ctx, userID, answers, userProfile, vehicles)
}

func (s *RecommendationService) generateWithAI(
	ctx context.Context,
	userID int64,
	answers []models.SubmittedAnswer,
	questions []models.Question,
	vehicles []models.Vehicle,
) (*models.Recommendation, error) {
	aiRec, err := s.ai.Recommend(ctx, answers, questions, vehicles)
	if err != nil {
		return nil, err
	}

	// Build vehicle lookup map.
	vehicleByID := make(map[int64]models.Vehicle, len(vehicles))
	for _, v := range vehicles {
		vehicleByID[v.ID] = v
	}

	items := make([]models.RecommendationItem, 0, len(aiRec.Items))
	for _, item := range aiRec.Items {
		v, ok := vehicleByID[item.VehicleID]
		if !ok {
			continue
		}
		// Keep the scoring data alongside AI reason.
		score, matches := scoreVehicle(
			map[string]float64{},
			v.MatchProfile,
		)
		_ = score
		items = append(items, models.RecommendationItem{
			Vehicle:         v,
			Rank:            item.Rank,
			Score:           0,
			Reason:          item.Reason,
			MatchedCriteria: matches,
		})
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("ai returned no recognizable vehicle IDs")
	}
	// Ensure rank is correct even if AI returned them out of order.
	sort.SliceStable(items, func(i, j int) bool { return items[i].Rank < items[j].Rank })
	for i := range items {
		items[i].Rank = i + 1
	}

	return &models.Recommendation{
		UserID:    userID,
		Summary:   "Recomendação gerada por Inteligência Artificial.",
		AISummary: aiRec.Summary,
		Items:     items,
	}, nil
}

func (s *RecommendationService) generateWithScoring(
	ctx context.Context,
	userID int64,
	answers []models.SubmittedAnswer,
	userProfile map[string]float64,
	vehicles []models.Vehicle,
) (*models.Recommendation, error) {
	items := make([]models.RecommendationItem, 0, len(vehicles))
	for _, vehicle := range vehicles {
		score, matches := scoreVehicle(userProfile, vehicle.MatchProfile)
		items = append(items, models.RecommendationItem{
			Vehicle:         vehicle,
			Score:           score,
			Reason:          buildReason(matches),
			MatchedCriteria: matches,
		})
	}

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Score == items[j].Score {
			return items[i].Vehicle.ID < items[j].Vehicle.ID
		}
		return items[i].Score > items[j].Score
	})
	if len(items) > 10 {
		items = items[:10]
	}
	for index := range items {
		items[index].Rank = index + 1
	}

	recommendation := &models.Recommendation{
		UserID:  userID,
		Summary: "Veículos ordenados pela compatibilidade com suas preferências.",
		Items:   items,
	}
	if err := s.recommendationRepo.Create(ctx, recommendation, answers); err != nil {
		return nil, err
	}
	return recommendation, nil
}

func (s *RecommendationService) History(ctx context.Context, userID int64, page, limit int) ([]models.Recommendation, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.recommendationRepo.GetByUser(ctx, userID, limit, (page-1)*limit)
}

func (s *RecommendationService) Get(ctx context.Context, id, userID int64) (*models.Recommendation, error) {
	return s.recommendationRepo.GetByID(ctx, id, userID)
}

func scoreVehicle(userProfile, vehicleProfile map[string]float64) (float64, []string) {
	var possible float64
	var matched float64
	matches := []string{}

	for dimension, userWeight := range userProfile {
		possible += userWeight
		vehicleWeight := vehicleProfile[dimension]
		if vehicleWeight <= 0 {
			continue
		}
		matched += userWeight * math.Min(vehicleWeight, 1)
		if vehicleWeight >= 0.6 {
			matches = append(matches, dimension)
		}
	}
	if possible == 0 {
		return 0, matches
	}
	return math.Round((matched/possible)*10000) / 100, matches
}

var dimensionLabels = map[string]string{
	"automatic":   "câmbio automático",
	"budget_high": "faixa de orçamento",
	"budget_low":  "faixa de orçamento",
	"budget_mid":  "faixa de orçamento",
	"cargo":       "capacidade de carga",
	"comfort":     "conforto",
	"compact":     "dimensões compactas",
	"efficiency":  "economia",
	"family":      "uso familiar",
	"manual":      "câmbio manual",
	"mixed":       "uso misto",
	"offroad":     "robustez",
	"performance": "desempenho",
	"space":       "espaço interno",
	"urban":       "uso urbano",
}

func buildReason(matches []string) string {
	if len(matches) == 0 {
		return "Compatibilidade geral com o perfil informado."
	}
	sort.Strings(matches)
	labels := make([]string, 0, len(matches))
	seen := map[string]bool{}
	for _, dimension := range matches {
		label := dimensionLabels[dimension]
		if label == "" || seen[label] {
			continue
		}
		seen[label] = true
		labels = append(labels, label)
	}
	if len(labels) > 4 {
		labels = labels[:4]
	}
	return fmt.Sprintf("Boa compatibilidade em %s.", strings.Join(labels, ", "))
}
