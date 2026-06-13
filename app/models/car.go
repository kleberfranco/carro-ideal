package models

import "time"

type VehicleCategory struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Active       bool   `json:"active"`
	VehicleCount int    `json:"vehicle_count,omitempty"`
}

type Vehicle struct {
	ID                  int64              `json:"id"`
	CategoryID          int64              `json:"category_id"`
	Category            VehicleCategory    `json:"category"`
	Brand               string             `json:"brand"`
	Model               string             `json:"model"`
	Version             string             `json:"version"`
	Year                int                `json:"year"`
	Condition           string             `json:"condition"`
	FuelType            string             `json:"fuel_type"`
	Transmission        string             `json:"transmission"`
	PriceMin            float64            `json:"price_min"`
	PriceMax            float64            `json:"price_max"`
	Seats               int                `json:"seats"`
	TrunkCapacity       int                `json:"trunk_capacity"`
	ConsumptionCity     float64            `json:"consumption_city"`
	ConsumptionHighway  float64            `json:"consumption_highway"`
	Description         string             `json:"description"`
	Strengths           string             `json:"strengths"`
	Weaknesses          string             `json:"weaknesses"`
	MatchProfile        map[string]float64 `json:"match_profile,omitempty"`
	Active              bool               `json:"active"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
	RecommendationCount int                `json:"recommendation_count,omitempty"`
}

type VehicleInput struct {
	CategoryID         *int64             `json:"category_id"`
	Brand              *string            `json:"brand"`
	Model              *string            `json:"model"`
	Version            *string            `json:"version"`
	Year               *int               `json:"year"`
	Condition          *string            `json:"condition"`
	FuelType           *string            `json:"fuel_type"`
	Transmission       *string            `json:"transmission"`
	PriceMin           *float64           `json:"price_min"`
	PriceMax           *float64           `json:"price_max"`
	Seats              *int               `json:"seats"`
	TrunkCapacity      *int               `json:"trunk_capacity"`
	ConsumptionCity    *float64           `json:"consumption_city"`
	ConsumptionHighway *float64           `json:"consumption_highway"`
	Description        *string            `json:"description"`
	Strengths          *string            `json:"strengths"`
	Weaknesses         *string            `json:"weaknesses"`
	MatchProfile       map[string]float64 `json:"match_profile"`
	Active             *bool              `json:"active"`
}

type CategoryInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Active      *bool   `json:"active"`
}
