package models

import "time"

type VehicleCategory struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Vehicle struct {
	ID                 int64              `json:"id"`
	CategoryID         int64              `json:"category_id"`
	Category           VehicleCategory    `json:"category"`
	Brand              string             `json:"brand"`
	Model              string             `json:"model"`
	Version            string             `json:"version"`
	Year               int                `json:"year"`
	FuelType           string             `json:"fuel_type"`
	Transmission       string             `json:"transmission"`
	PriceMin           float64            `json:"price_min"`
	PriceMax           float64            `json:"price_max"`
	Seats              int                `json:"seats"`
	TrunkCapacity      int                `json:"trunk_capacity"`
	ConsumptionCity    float64            `json:"consumption_city"`
	ConsumptionHighway float64            `json:"consumption_highway"`
	Description        string             `json:"description"`
	Strengths          string             `json:"strengths"`
	Weaknesses         string             `json:"weaknesses"`
	MatchProfile       map[string]float64 `json:"-"`
	Active             bool               `json:"active"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
}
