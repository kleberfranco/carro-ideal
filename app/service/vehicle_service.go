package service

import (
	"context"

	"carro-ideal/app/models"
	"carro-ideal/app/repository"
)

type VehicleService struct {
	repo  repository.VehicleRepository
	cache *CatalogCache
}

func NewVehicleService(repo repository.VehicleRepository, caches ...*CatalogCache) *VehicleService {
	var cache *CatalogCache
	if len(caches) > 0 {
		cache = caches[0]
	}
	return &VehicleService{repo: repo, cache: cache}
}

func (s *VehicleService) GetActive(ctx context.Context, categoryID int64) ([]models.Vehicle, error) {
	if s.cache != nil {
		if vehicles, found := s.cache.Vehicles(categoryID); found {
			return vehicles, nil
		}
	}
	vehicles, err := s.repo.GetActive(ctx, categoryID)
	if err == nil && s.cache != nil {
		s.cache.SetVehicles(categoryID, vehicles)
	}
	return vehicles, err
}

func (s *VehicleService) GetByID(ctx context.Context, id int64) (*models.Vehicle, error) {
	return s.repo.GetByID(ctx, id)
}
