package service

import (
	"context"

	"carro-ideal/app/models"
	"carro-ideal/app/repository"
)

type VehicleService struct {
	repo repository.VehicleRepository
}

func NewVehicleService(repo repository.VehicleRepository) *VehicleService {
	return &VehicleService{repo: repo}
}

func (s *VehicleService) GetActive(ctx context.Context, categoryID int64) ([]models.Vehicle, error) {
	return s.repo.GetActive(ctx, categoryID)
}

func (s *VehicleService) GetByID(ctx context.Context, id int64) (*models.Vehicle, error) {
	return s.repo.GetByID(ctx, id)
}
