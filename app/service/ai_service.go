package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"carro-ideal/app/models"
)

// AIClient is the interface satisfied by clients.OpenAIClient (and test fakes).
type AIClient interface {
	ChatComplete(ctx context.Context, systemPrompt, userPrompt string) (string, error)
}

// AIRecommendationItem holds one vehicle ranked by ChatGPT.
type AIRecommendationItem struct {
	VehicleID int64  `json:"vehicle_id"`
	Rank      int    `json:"rank"`
	Reason    string `json:"reason"`
}

// AIRecommendation is the structured response expected from ChatGPT.
type AIRecommendation struct {
	Summary string                 `json:"summary"`
	Items   []AIRecommendationItem `json:"items"`
}

// AIService uses the ChatGPT API to generate vehicle recommendations.
type AIService struct {
	client AIClient
}

func NewAIService(client AIClient) *AIService {
	return &AIService{client: client}
}

const systemPrompt = `Você é um especialista em automóveis brasileiro com amplo conhecimento do mercado nacional.
Sua tarefa é recomendar os veículos mais adequados ao perfil do usuário com base nas respostas do questionário e no catálogo disponível.
Responda APENAS com JSON válido, sem texto adicional antes ou depois.`

// Recommend calls ChatGPT with the user's answers and vehicle catalog, returning a ranked AIRecommendation.
func (s *AIService) Recommend(
	ctx context.Context,
	answers []models.SubmittedAnswer,
	questions []models.Question,
	vehicles []models.Vehicle,
) (*AIRecommendation, error) {
	userPrompt := buildPrompt(answers, questions, vehicles)
	raw, err := s.client.ChatComplete(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("chatgpt: %w", err)
	}
	return parseAIResponse(raw)
}

func buildPrompt(answers []models.SubmittedAnswer, questions []models.Question, vehicles []models.Vehicle) string {
	var sb strings.Builder

	sb.WriteString("Respostas do usuário ao questionário:\n")
	answerMap := make(map[int64]int64, len(answers))
	for _, a := range answers {
		answerMap[a.QuestionID] = a.AnswerOptionID
	}
	for _, q := range questions {
		selectedID := answerMap[q.ID]
		for _, opt := range q.Options {
			if opt.ID == selectedID {
				sb.WriteString(fmt.Sprintf("- %s: %s\n", q.Text, opt.Text))
				break
			}
		}
	}

	sb.WriteString("\nCatálogo de veículos disponíveis:\n")
	for _, v := range vehicles {
		price := ""
		if v.PriceMin > 0 && v.PriceMax > 0 {
			price = fmt.Sprintf(" — R$ %.0f a R$ %.0f", v.PriceMin, v.PriceMax)
		} else if v.PriceMin > 0 {
			price = fmt.Sprintf(" — a partir de R$ %.0f", v.PriceMin)
		}
		transmission := map[string]string{
			"automatic": "automático",
			"manual":    "manual",
			"cvt":       "CVT",
		}[v.Transmission]
		if transmission == "" {
			transmission = v.Transmission
		}
		line := fmt.Sprintf("[ID:%d] %s %s %d — %s%s", v.ID, v.Brand, v.Model, v.Year, v.Category.Name, price)
		if transmission != "" {
			line += ", câmbio " + transmission
		}
		sb.WriteString(line + "\n")
		if v.Description != "" {
			sb.WriteString("  " + v.Description + "\n")
		}
		if v.Strengths != "" {
			sb.WriteString("  Pontos fortes: " + v.Strengths + "\n")
		}
	}

	sb.WriteString(`
Responda com o seguinte JSON (sem nenhum texto fora do JSON):
{
  "summary": "parágrafo em português descrevendo o perfil do usuário e por que estes veículos foram escolhidos",
  "items": [
    {"vehicle_id": <id>, "rank": 1, "reason": "justificativa em português para este veículo"},
    {"vehicle_id": <id>, "rank": 2, "reason": "..."},
    {"vehicle_id": <id>, "rank": 3, "reason": "..."}
  ]
}
Inclua no máximo 5 veículos. Ordene do mais ao menos adequado. Use apenas IDs dos veículos listados acima.`)

	return sb.String()
}

func parseAIResponse(raw string) (*AIRecommendation, error) {
	raw = strings.TrimSpace(raw)
	// strip markdown code fences if present
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	var rec AIRecommendation
	if err := json.Unmarshal([]byte(raw), &rec); err != nil {
		return nil, fmt.Errorf("parse ai response: %w", err)
	}
	if len(rec.Items) == 0 {
		return nil, fmt.Errorf("ai returned no recommendation items")
	}
	return &rec, nil
}
