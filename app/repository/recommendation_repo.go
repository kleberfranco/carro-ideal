package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"carro-ideal/app/models"
)

var ErrRecommendationNotFound = errors.New("recommendation not found")

type RecommendationRepository interface {
	Create(ctx context.Context, recommendation *models.Recommendation, answers []models.SubmittedAnswer) error
	GetByUser(ctx context.Context, userID int64, limit, offset int) ([]models.Recommendation, int, error)
	GetByID(ctx context.Context, id, userID int64) (*models.Recommendation, error)
}

type recommendationRepository struct {
	db *sql.DB
}

func NewRecommendationRepository(db *sql.DB) RecommendationRepository {
	return &recommendationRepository{db: db}
}

func (r *recommendationRepository) Create(ctx context.Context, recommendation *models.Recommendation, answers []models.SubmittedAnswer) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	snapshot, err := json.Marshal(answers)
	if err != nil {
		return err
	}
	err = tx.QueryRowContext(ctx, `
		INSERT INTO recommendations (user_id, answers_snapshot, summary, ai_summary)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`,
		recommendation.UserID, snapshot, recommendation.Summary, recommendation.AISummary,
	).Scan(&recommendation.ID, &recommendation.CreatedAt)
	if err != nil {
		return err
	}

	for _, answer := range answers {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO user_answers
			    (user_id, recommendation_id, question_id, answer_option_id)
			VALUES ($1, $2, $3, $4)`,
			recommendation.UserID, recommendation.ID, answer.QuestionID, answer.AnswerOptionID,
		); err != nil {
			return err
		}
	}

	for index := range recommendation.Items {
		item := &recommendation.Items[index]
		item.RecommendationID = recommendation.ID
		err := tx.QueryRowContext(ctx, `
			INSERT INTO recommendation_items
			    (recommendation_id, vehicle_id, rank, score, reason)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at`,
			recommendation.ID, item.Vehicle.ID, item.Rank, item.Score, item.Reason,
		).Scan(&item.ID, &item.CreatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *recommendationRepository) GetByUser(ctx context.Context, userID int64, limit, offset int) ([]models.Recommendation, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM recommendations WHERE user_id=$1", userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT r.id, r.user_id, r.summary, r.ai_summary, r.created_at, COUNT(i.id)
		FROM recommendations r
		LEFT JOIN recommendation_items i ON i.recommendation_id=r.id
		WHERE r.user_id=$1
		GROUP BY r.id
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	recommendations := []models.Recommendation{}
	for rows.Next() {
		var recommendation models.Recommendation
		if err := rows.Scan(
			&recommendation.ID,
			&recommendation.UserID,
			&recommendation.Summary,
			&recommendation.AISummary,
			&recommendation.CreatedAt,
			&recommendation.ItemCount,
		); err != nil {
			return nil, 0, err
		}
		recommendations = append(recommendations, recommendation)
	}
	return recommendations, total, rows.Err()
}

func (r *recommendationRepository) GetByID(ctx context.Context, id, userID int64) (*models.Recommendation, error) {
	recommendation := &models.Recommendation{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, summary, ai_summary, created_at
		FROM recommendations WHERE id=$1 AND user_id=$2`, id, userID,
	).Scan(&recommendation.ID, &recommendation.UserID, &recommendation.Summary, &recommendation.AISummary, &recommendation.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRecommendationNotFound
	}
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT i.id, i.recommendation_id, i.rank, i.score, i.reason, i.created_at,
		       `+vehicleColumns+`
		FROM recommendation_items i
		JOIN vehicles v ON v.id=i.vehicle_id
		JOIN vehicle_categories c ON c.id=v.category_id
		WHERE i.recommendation_id=$1
		ORDER BY i.rank`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recommendation.Items = []models.RecommendationItem{}
	for rows.Next() {
		var item models.RecommendationItem
		var profile []byte
		err := rows.Scan(
			&item.ID, &item.RecommendationID, &item.Rank, &item.Score, &item.Reason, &item.CreatedAt,
			&item.Vehicle.ID, &item.Vehicle.CategoryID, &item.Vehicle.Category.ID,
			&item.Vehicle.Category.Name, &item.Vehicle.Category.Description, &item.Vehicle.Category.Active, &item.Vehicle.Brand,
			&item.Vehicle.Model, &item.Vehicle.Version, &item.Vehicle.Year, &item.Vehicle.FuelType,
			&item.Vehicle.Transmission, &item.Vehicle.PriceMin, &item.Vehicle.PriceMax,
			&item.Vehicle.Seats, &item.Vehicle.TrunkCapacity, &item.Vehicle.ConsumptionCity,
			&item.Vehicle.ConsumptionHighway, &item.Vehicle.Description, &item.Vehicle.Strengths,
			&item.Vehicle.Weaknesses, &profile, &item.Vehicle.Active, &item.Vehicle.CreatedAt,
			&item.Vehicle.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(profile, &item.Vehicle.MatchProfile); err != nil {
			return nil, err
		}
		recommendation.Items = append(recommendation.Items, item)
	}
	return recommendation, rows.Err()
}
