package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"carro-ideal/app/models"
)

var ErrVehicleNotFound = errors.New("vehicle not found")

type VehicleRepository interface {
	GetActive(ctx context.Context, categoryID int64) ([]models.Vehicle, error)
	GetByID(ctx context.Context, id int64) (*models.Vehicle, error)
}

type vehicleRepository struct {
	db *sql.DB
}

func NewVehicleRepository(db *sql.DB) VehicleRepository {
	return &vehicleRepository{db: db}
}

const vehicleColumns = `
	v.id, v.category_id, c.id, c.name, COALESCE(c.description, ''),
	v.brand, v.model, COALESCE(v.version, ''), v.year, v.fuel_type, v.transmission,
	v.price_min, v.price_max, COALESCE(v.seats, 0), COALESCE(v.trunk_capacity, 0),
	COALESCE(v.consumption_city, 0), COALESCE(v.consumption_highway, 0),
	COALESCE(v.description, ''), COALESCE(v.strengths, ''), COALESCE(v.weaknesses, ''),
	v.match_profile, v.active, v.created_at, v.updated_at`

func (r *vehicleRepository) GetActive(ctx context.Context, categoryID int64) ([]models.Vehicle, error) {
	query := `SELECT ` + vehicleColumns + `
		FROM vehicles v
		JOIN vehicle_categories c ON c.id = v.category_id
		WHERE v.active=true AND c.active=true`
	args := []interface{}{}
	if categoryID > 0 {
		query += " AND v.category_id=$1"
		args = append(args, categoryID)
	}
	query += " ORDER BY v.brand, v.model"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vehicles := []models.Vehicle{}
	for rows.Next() {
		vehicle, err := scanVehicle(rows)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, *vehicle)
	}
	return vehicles, rows.Err()
}

func (r *vehicleRepository) GetByID(ctx context.Context, id int64) (*models.Vehicle, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+vehicleColumns+`
		FROM vehicles v
		JOIN vehicle_categories c ON c.id = v.category_id
		WHERE v.id=$1 AND v.active=true AND c.active=true`, id)
	vehicle, err := scanVehicle(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrVehicleNotFound
	}
	return vehicle, err
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanVehicle(scanner rowScanner) (*models.Vehicle, error) {
	vehicle := &models.Vehicle{}
	var profile []byte
	err := scanner.Scan(
		&vehicle.ID, &vehicle.CategoryID, &vehicle.Category.ID, &vehicle.Category.Name,
		&vehicle.Category.Description, &vehicle.Brand, &vehicle.Model, &vehicle.Version,
		&vehicle.Year, &vehicle.FuelType, &vehicle.Transmission, &vehicle.PriceMin,
		&vehicle.PriceMax, &vehicle.Seats, &vehicle.TrunkCapacity, &vehicle.ConsumptionCity,
		&vehicle.ConsumptionHighway, &vehicle.Description, &vehicle.Strengths,
		&vehicle.Weaknesses, &profile, &vehicle.Active, &vehicle.CreatedAt, &vehicle.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(profile, &vehicle.MatchProfile); err != nil {
		return nil, err
	}
	return vehicle, nil
}
