package health

import (
	"database/sql"
	"net/http"
	"time"

	"carro-ideal/app/internal/response"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(database *sql.DB) *Handler {
	return &Handler{db: database}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if h.db == nil || h.db.PingContext(r.Context()) != nil {
		response.Error(w, http.StatusServiceUnavailable, "database unavailable", "DB_UNAVAILABLE")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
