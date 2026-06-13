package models

import "time"

type Question struct {
	ID           int64          `json:"id"`
	Text         string         `json:"text"`
	Type         string         `json:"type"`
	Weight       float64        `json:"weight"`
	DisplayOrder int            `json:"display_order"`
	Active       bool           `json:"active"`
	Options      []AnswerOption `json:"answer_options"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type AnswerOption struct {
	ID           int64              `json:"id"`
	QuestionID   int64              `json:"question_id"`
	Text         string             `json:"text"`
	ScoreProfile map[string]float64 `json:"score_profile"`
	DisplayOrder int                `json:"display_order"`
	Active       bool               `json:"active"`
}

type QuestionInput struct {
	Text         *string  `json:"text"`
	Type         *string  `json:"type"`
	Weight       *float64 `json:"weight"`
	DisplayOrder *int     `json:"display_order"`
	Active       *bool    `json:"active"`
}

type AnswerOptionInput struct {
	Text         *string            `json:"text"`
	ScoreProfile map[string]float64 `json:"score_profile"`
	DisplayOrder *int               `json:"display_order"`
	Active       *bool              `json:"active"`
}

type SubmittedAnswer struct {
	QuestionID     int64 `json:"question_id"`
	AnswerOptionID int64 `json:"answer_option_id"`
}
