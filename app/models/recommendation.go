package models

import "time"

type Recommendation struct {
	ID        int64                `json:"id"`
	UserID    int64                `json:"user_id"`
	Summary   string               `json:"summary"`
	Items     []RecommendationItem `json:"items,omitempty"`
	ItemCount int                  `json:"item_count,omitempty"`
	CreatedAt time.Time            `json:"created_at"`
}

type RecommendationItem struct {
	ID               int64     `json:"id"`
	RecommendationID int64     `json:"recommendation_id"`
	Vehicle          Vehicle   `json:"vehicle"`
	Rank             int       `json:"rank"`
	Score            float64   `json:"score"`
	Reason           string    `json:"reason"`
	MatchedCriteria  []string  `json:"matched_criteria,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}
