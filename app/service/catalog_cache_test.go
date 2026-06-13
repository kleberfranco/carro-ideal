package service

import (
	"testing"
	"time"

	"carro-ideal/app/models"
)

func TestCatalogCacheInvalidation(t *testing.T) {
	cache := NewCatalogCache(time.Minute)
	cache.SetVehicles(0, []models.Vehicle{{ID: 1}})

	if _, found := cache.Vehicles(0); !found {
		t.Fatal("vehicle cache should contain the stored value")
	}
	cache.Invalidate()
	if _, found := cache.Vehicles(0); found {
		t.Fatal("vehicle cache should be empty after invalidation")
	}
}

func TestCatalogCacheQuestions(t *testing.T) {
	cache := NewCatalogCache(time.Minute)

	if _, found := cache.Questions(); found {
		t.Fatal("cache should be empty initially")
	}

	questions := []models.Question{{ID: 1, Text: "Qual o seu orçamento?"}}
	cache.SetQuestions(questions)

	got, found := cache.Questions()
	if !found {
		t.Fatal("cache should return questions after SetQuestions")
	}
	if len(got) != 1 || got[0].ID != 1 {
		t.Fatalf("cache returned wrong questions: %v", got)
	}

	cache.Invalidate()
	if _, found := cache.Questions(); found {
		t.Fatal("questions cache should be empty after invalidation")
	}
}

func TestCatalogCacheExpiredTTL(t *testing.T) {
	cache := NewCatalogCache(0)
	cache.SetVehicles(0, []models.Vehicle{{ID: 1}})
	cache.SetQuestions([]models.Question{{ID: 1}})

	if _, found := cache.Vehicles(0); found {
		t.Fatal("zero-TTL cache should always report miss for vehicles")
	}
	if _, found := cache.Questions(); found {
		t.Fatal("zero-TTL cache should always report miss for questions")
	}
}
