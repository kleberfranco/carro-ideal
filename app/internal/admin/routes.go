package admin

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Get("/dashboard", h.Dashboard)
	r.Get("/users", h.Users)
	r.Get("/vehicles", h.Vehicles)
	r.Post("/vehicles", h.CreateVehicle)
	r.Put("/vehicles/{id}", h.UpdateVehicle)
	r.Delete("/vehicles/{id}", h.DeleteVehicle)
	r.Get("/categories", h.Categories)
	r.Post("/categories", h.CreateCategory)
	r.Put("/categories/{id}", h.UpdateCategory)
	r.Delete("/categories/{id}", h.DeleteCategory)
	r.Get("/questions", h.Questions)
	r.Post("/questions", h.CreateQuestion)
	r.Put("/questions/{id}", h.UpdateQuestion)
	r.Delete("/questions/{id}", h.DeleteQuestion)
	r.Post("/questions/{id}/options", h.CreateOption)
	r.Put("/questions/{id}/options/{optionId}", h.UpdateOption)
	r.Delete("/questions/{id}/options/{optionId}", h.DeleteOption)
}
