package health

import (
	"net/http"
	"time"

	"carro-ideal/app/db"
	"carro-ideal/app/internal/response"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if err := db.GetDBHealth(r.Context()); err != nil {
		response.Error(w, http.StatusServiceUnavailable, "database unavailable", "DB_UNAVAILABLE")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
