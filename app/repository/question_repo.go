package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"carro-ideal/app/models"
)

type QuestionRepository interface {
	GetActive(ctx context.Context) ([]models.Question, error)
}

type questionRepository struct {
	db *sql.DB
}

func NewQuestionRepository(db *sql.DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) GetActive(ctx context.Context) ([]models.Question, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT q.id, q.question_text, q.question_type, q.weight, q.display_order,
		       q.active, q.created_at, q.updated_at,
		       o.id, o.option_text, o.score_profile, o.display_order, o.active
		FROM questions q
		JOIN answer_options o ON o.question_id=q.id AND o.active=true
		WHERE q.active=true
		ORDER BY q.display_order, o.display_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	questions := []models.Question{}
	index := map[int64]int{}
	for rows.Next() {
		var question models.Question
		var option models.AnswerOption
		var profile []byte
		if err := rows.Scan(
			&question.ID, &question.Text, &question.Type, &question.Weight,
			&question.DisplayOrder, &question.Active, &question.CreatedAt, &question.UpdatedAt,
			&option.ID, &option.Text, &profile, &option.DisplayOrder, &option.Active,
		); err != nil {
			return nil, err
		}
		option.QuestionID = question.ID
		if err := json.Unmarshal(profile, &option.ScoreProfile); err != nil {
			return nil, err
		}

		position, exists := index[question.ID]
		if !exists {
			question.Options = []models.AnswerOption{}
			questions = append(questions, question)
			position = len(questions) - 1
			index[question.ID] = position
		}
		questions[position].Options = append(questions[position].Options, option)
	}
	return questions, rows.Err()
}
