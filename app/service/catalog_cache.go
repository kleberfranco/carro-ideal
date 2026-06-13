package service

import (
	"sync"
	"time"

	"carro-ideal/app/models"
)

type CatalogCache struct {
	mu              sync.RWMutex
	ttl             time.Duration
	questions       []models.Question
	questionsExpiry time.Time
	vehicles        map[int64]vehicleCacheEntry
}

type vehicleCacheEntry struct {
	items     []models.Vehicle
	expiresAt time.Time
}

func NewCatalogCache(ttl time.Duration) *CatalogCache {
	return &CatalogCache{ttl: ttl, vehicles: map[int64]vehicleCacheEntry{}}
}

func (c *CatalogCache) Questions() ([]models.Question, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.expired(c.questionsExpiry) || c.questions == nil {
		return nil, false
	}
	return c.questions, true
}

func (c *CatalogCache) SetQuestions(questions []models.Question) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.questions = questions
	c.questionsExpiry = time.Now().Add(c.ttl)
}

func (c *CatalogCache) Vehicles(categoryID int64) ([]models.Vehicle, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, found := c.vehicles[categoryID]
	if !found || c.expired(entry.expiresAt) {
		return nil, false
	}
	return entry.items, true
}

func (c *CatalogCache) SetVehicles(categoryID int64, vehicles []models.Vehicle) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.vehicles[categoryID] = vehicleCacheEntry{
		items:     vehicles,
		expiresAt: time.Now().Add(c.ttl),
	}
}

func (c *CatalogCache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.questions = nil
	c.questionsExpiry = time.Time{}
	c.vehicles = map[int64]vehicleCacheEntry{}
}

func (c *CatalogCache) expired(expiresAt time.Time) bool {
	return c.ttl <= 0 || time.Now().After(expiresAt)
}
