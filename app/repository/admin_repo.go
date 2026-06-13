package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"carro-ideal/app/models"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryInUse    = errors.New("category in use")
	ErrQuestionNotFound = errors.New("question not found")
	ErrOptionNotFound   = errors.New("answer option not found")
)

type AdminRepositoryInterface interface {
	Stats(ctx context.Context) (*models.AdminStats, error)
	Vehicles(ctx context.Context, search string, limit, offset int) ([]models.Vehicle, int, error)
	Vehicle(ctx context.Context, id int64) (*models.Vehicle, error)
	CreateVehicle(ctx context.Context, vehicle *models.Vehicle) error
	UpdateVehicle(ctx context.Context, vehicle *models.Vehicle) error
	DeleteVehicle(ctx context.Context, id int64) error
	Categories(ctx context.Context) ([]models.VehicleCategory, error)
	Category(ctx context.Context, id int64) (*models.VehicleCategory, error)
	CreateCategory(ctx context.Context, category *models.VehicleCategory) error
	UpdateCategory(ctx context.Context, category *models.VehicleCategory) error
	DeleteCategory(ctx context.Context, id int64) error
	Questions(ctx context.Context) ([]models.Question, error)
	Question(ctx context.Context, id int64) (*models.Question, error)
	CreateQuestion(ctx context.Context, question *models.Question) error
	UpdateQuestion(ctx context.Context, question *models.Question) error
	DeleteQuestion(ctx context.Context, id int64) error
	CreateOption(ctx context.Context, option *models.AnswerOption) error
	UpdateOption(ctx context.Context, option *models.AnswerOption) error
	DeleteOption(ctx context.Context, questionID, optionID int64) error
}

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) AdminRepositoryInterface {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) Stats(ctx context.Context) (*models.AdminStats, error) {
	stats := &models.AdminStats{}
	err := r.db.QueryRowContext(ctx, `
		SELECT
		    (SELECT COUNT(*) FROM users),
		    (SELECT COUNT(*) FROM vehicles),
		    (SELECT COUNT(*) FROM recommendations),
		    (SELECT COUNT(*) FROM questions),
		    (SELECT COUNT(DISTINCT user_id) FROM sessions WHERE created_at >= NOW() - INTERVAL '7 days'),
		    (SELECT COUNT(*) FROM users WHERE created_at >= NOW() - INTERVAL '7 days')`,
	).Scan(
		&stats.Users, &stats.Vehicles, &stats.Recommendations, &stats.Questions,
		&stats.ActiveUsersWeek, &stats.NewUsersWeek,
	)
	return stats, err
}

func (r *AdminRepository) Vehicles(ctx context.Context, search string, limit, offset int) ([]models.Vehicle, int, error) {
	pattern := "%" + search + "%"
	var total int
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM vehicles
		WHERE $1='' OR brand ILIKE $2 OR model ILIKE $2`, search, pattern,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx, `SELECT `+vehicleColumns+`,
		(SELECT COUNT(*) FROM recommendation_items i WHERE i.vehicle_id=v.id)
		FROM vehicles v
		JOIN vehicle_categories c ON c.id=v.category_id
		WHERE $1='' OR v.brand ILIKE $2 OR v.model ILIKE $2
		ORDER BY v.active DESC, v.brand, v.model
		LIMIT $3 OFFSET $4`, search, pattern, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	vehicles := []models.Vehicle{}
	for rows.Next() {
		vehicle := &models.Vehicle{}
		var profile []byte
		err := rows.Scan(
			&vehicle.ID, &vehicle.CategoryID, &vehicle.Category.ID, &vehicle.Category.Name,
			&vehicle.Category.Description, &vehicle.Category.Active, &vehicle.Brand, &vehicle.Model, &vehicle.Version,
			&vehicle.Year, &vehicle.Condition, &vehicle.FuelType, &vehicle.Transmission, &vehicle.PriceMin,
			&vehicle.PriceMax, &vehicle.Seats, &vehicle.TrunkCapacity, &vehicle.ConsumptionCity,
			&vehicle.ConsumptionHighway, &vehicle.Description, &vehicle.Strengths,
			&vehicle.Weaknesses, &profile, &vehicle.Active, &vehicle.CreatedAt, &vehicle.UpdatedAt,
			&vehicle.RecommendationCount,
		)
		if err != nil {
			return nil, 0, err
		}
		if err := json.Unmarshal(profile, &vehicle.MatchProfile); err != nil {
			return nil, 0, err
		}
		vehicles = append(vehicles, *vehicle)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("list admin vehicles: %w", err)
	}
	return vehicles, total, nil
}

func (r *AdminRepository) Vehicle(ctx context.Context, id int64) (*models.Vehicle, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+vehicleColumns+`
		FROM vehicles v JOIN vehicle_categories c ON c.id=v.category_id WHERE v.id=$1`, id)
	vehicle, err := scanVehicle(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrVehicleNotFound
	}
	return vehicle, err
}

func (r *AdminRepository) CreateVehicle(ctx context.Context, vehicle *models.Vehicle) error {
	profile, err := json.Marshal(vehicle.MatchProfile)
	if err != nil {
		return err
	}
	return r.db.QueryRowContext(ctx, `
		INSERT INTO vehicles (
		    category_id, brand, model, version, year, fuel_type, transmission,
		    price_min, price_max, seats, trunk_capacity, consumption_city,
		    consumption_highway, description, strengths, weaknesses, match_profile, active
		) VALUES (
		    $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18
		) RETURNING id, created_at, updated_at`,
		vehicle.CategoryID, vehicle.Brand, vehicle.Model, vehicle.Version, vehicle.Year,
		vehicle.FuelType, vehicle.Transmission, vehicle.PriceMin, vehicle.PriceMax,
		vehicle.Seats, vehicle.TrunkCapacity, vehicle.ConsumptionCity,
		vehicle.ConsumptionHighway, vehicle.Description, vehicle.Strengths,
		vehicle.Weaknesses, profile, vehicle.Active,
	).Scan(&vehicle.ID, &vehicle.CreatedAt, &vehicle.UpdatedAt)
}

func (r *AdminRepository) UpdateVehicle(ctx context.Context, vehicle *models.Vehicle) error {
	profile, err := json.Marshal(vehicle.MatchProfile)
	if err != nil {
		return err
	}
	result, err := r.db.ExecContext(ctx, `
		UPDATE vehicles SET
		    category_id=$1, brand=$2, model=$3, version=$4, year=$5,
		    fuel_type=$6, transmission=$7, price_min=$8, price_max=$9,
		    seats=$10, trunk_capacity=$11, consumption_city=$12,
		    consumption_highway=$13, description=$14, strengths=$15,
		    weaknesses=$16, match_profile=$17, active=$18, updated_at=NOW()
		WHERE id=$19`,
		vehicle.CategoryID, vehicle.Brand, vehicle.Model, vehicle.Version, vehicle.Year,
		vehicle.FuelType, vehicle.Transmission, vehicle.PriceMin, vehicle.PriceMax,
		vehicle.Seats, vehicle.TrunkCapacity, vehicle.ConsumptionCity,
		vehicle.ConsumptionHighway, vehicle.Description, vehicle.Strengths,
		vehicle.Weaknesses, profile, vehicle.Active, vehicle.ID,
	)
	return requireAffected(result, err, ErrVehicleNotFound)
}

func (r *AdminRepository) DeleteVehicle(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "UPDATE vehicles SET active=false, updated_at=NOW() WHERE id=$1", id)
	return requireAffected(result, err, ErrVehicleNotFound)
}

func (r *AdminRepository) Categories(ctx context.Context) ([]models.VehicleCategory, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT c.id, c.name, COALESCE(c.description,''), c.active, COUNT(v.id)
		FROM vehicle_categories c
		LEFT JOIN vehicles v ON v.category_id=c.id AND v.active=true
		GROUP BY c.id
		ORDER BY c.active DESC, c.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []models.VehicleCategory{}
	for rows.Next() {
		var category models.VehicleCategory
		if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.Active, &category.VehicleCount); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func (r *AdminRepository) Category(ctx context.Context, id int64) (*models.VehicleCategory, error) {
	category := &models.VehicleCategory{}
	err := r.db.QueryRowContext(ctx, `
		SELECT c.id, c.name, COALESCE(c.description,''), c.active,
		       (SELECT COUNT(*) FROM vehicles v WHERE v.category_id=c.id AND v.active=true)
		FROM vehicle_categories c WHERE c.id=$1`, id,
	).Scan(&category.ID, &category.Name, &category.Description, &category.Active, &category.VehicleCount)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrCategoryNotFound
	}
	return category, err
}

func (r *AdminRepository) CreateCategory(ctx context.Context, category *models.VehicleCategory) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO vehicle_categories (name, description, active)
		VALUES ($1,$2,$3) RETURNING id`,
		category.Name, category.Description, category.Active,
	).Scan(&category.ID)
}

func (r *AdminRepository) UpdateCategory(ctx context.Context, category *models.VehicleCategory) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE vehicle_categories
		SET name=$1, description=$2, active=$3, updated_at=NOW()
		WHERE id=$4`, category.Name, category.Description, category.Active, category.ID)
	return requireAffected(result, err, ErrCategoryNotFound)
}

func (r *AdminRepository) DeleteCategory(ctx context.Context, id int64) error {
	var count int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM vehicles WHERE category_id=$1 AND active=true", id).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return ErrCategoryInUse
	}
	result, err := r.db.ExecContext(ctx, "UPDATE vehicle_categories SET active=false, updated_at=NOW() WHERE id=$1", id)
	return requireAffected(result, err, ErrCategoryNotFound)
}

func (r *AdminRepository) Questions(ctx context.Context) ([]models.Question, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT q.id, q.question_text, q.question_type, q.weight, q.display_order,
		       q.active, q.created_at, q.updated_at,
		       o.id, COALESCE(o.option_text,''), COALESCE(o.score_profile,'{}'::jsonb),
		       COALESCE(o.display_order,0), COALESCE(o.active,false)
		FROM questions q
		LEFT JOIN answer_options o ON o.question_id=q.id
		ORDER BY q.display_order, o.display_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collectQuestions(rows)
}

func (r *AdminRepository) Question(ctx context.Context, id int64) (*models.Question, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT q.id, q.question_text, q.question_type, q.weight, q.display_order,
		       q.active, q.created_at, q.updated_at,
		       o.id, COALESCE(o.option_text,''), COALESCE(o.score_profile,'{}'::jsonb),
		       COALESCE(o.display_order,0), COALESCE(o.active,false)
		FROM questions q
		LEFT JOIN answer_options o ON o.question_id=q.id
		WHERE q.id=$1 ORDER BY o.display_order`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	questions, err := collectQuestions(rows)
	if err != nil {
		return nil, err
	}
	if len(questions) == 0 {
		return nil, ErrQuestionNotFound
	}
	return &questions[0], nil
}

func (r *AdminRepository) CreateQuestion(ctx context.Context, question *models.Question) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO questions (question_text, question_type, weight, display_order, active)
		VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at, updated_at`,
		question.Text, question.Type, question.Weight, question.DisplayOrder, question.Active,
	).Scan(&question.ID, &question.CreatedAt, &question.UpdatedAt)
}

func (r *AdminRepository) UpdateQuestion(ctx context.Context, question *models.Question) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE questions SET question_text=$1, question_type=$2, weight=$3,
		    display_order=$4, active=$5, updated_at=NOW() WHERE id=$6`,
		question.Text, question.Type, question.Weight, question.DisplayOrder, question.Active, question.ID,
	)
	return requireAffected(result, err, ErrQuestionNotFound)
}

func (r *AdminRepository) DeleteQuestion(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "UPDATE questions SET active=false, updated_at=NOW() WHERE id=$1", id)
	return requireAffected(result, err, ErrQuestionNotFound)
}

func (r *AdminRepository) CreateOption(ctx context.Context, option *models.AnswerOption) error {
	profile, err := json.Marshal(option.ScoreProfile)
	if err != nil {
		return err
	}
	return r.db.QueryRowContext(ctx, `
		INSERT INTO answer_options
		    (question_id, option_text, score_profile, display_order, active)
		VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		option.QuestionID, option.Text, profile, option.DisplayOrder, option.Active,
	).Scan(&option.ID)
}

func (r *AdminRepository) UpdateOption(ctx context.Context, option *models.AnswerOption) error {
	profile, err := json.Marshal(option.ScoreProfile)
	if err != nil {
		return err
	}
	result, err := r.db.ExecContext(ctx, `
		UPDATE answer_options SET option_text=$1, score_profile=$2,
		    display_order=$3, active=$4, updated_at=NOW()
		WHERE id=$5 AND question_id=$6`,
		option.Text, profile, option.DisplayOrder, option.Active, option.ID, option.QuestionID,
	)
	return requireAffected(result, err, ErrOptionNotFound)
}

func (r *AdminRepository) DeleteOption(ctx context.Context, questionID, optionID int64) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE answer_options SET active=false, updated_at=NOW()
		WHERE id=$1 AND question_id=$2`, optionID, questionID)
	return requireAffected(result, err, ErrOptionNotFound)
}

func collectQuestions(rows *sql.Rows) ([]models.Question, error) {
	questions := []models.Question{}
	index := map[int64]int{}
	for rows.Next() {
		var question models.Question
		var optionID sql.NullInt64
		var optionText string
		var profile []byte
		var optionOrder int
		var optionActive bool
		if err := rows.Scan(
			&question.ID, &question.Text, &question.Type, &question.Weight,
			&question.DisplayOrder, &question.Active, &question.CreatedAt, &question.UpdatedAt,
			&optionID, &optionText, &profile, &optionOrder, &optionActive,
		); err != nil {
			return nil, err
		}
		position, exists := index[question.ID]
		if !exists {
			question.Options = []models.AnswerOption{}
			questions = append(questions, question)
			position = len(questions) - 1
			index[question.ID] = position
		}
		if optionID.Valid {
			option := models.AnswerOption{
				ID:           optionID.Int64,
				QuestionID:   question.ID,
				Text:         optionText,
				DisplayOrder: optionOrder,
				Active:       optionActive,
			}
			if err := json.Unmarshal(profile, &option.ScoreProfile); err != nil {
				return nil, err
			}
			questions[position].Options = append(questions[position].Options, option)
		}
	}
	return questions, rows.Err()
}

func requireAffected(result sql.Result, err error, notFound error) error {
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return notFound
	}
	return nil
}
