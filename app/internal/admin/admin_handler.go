package admin

import (
	"net/http"
)

func (h *Handler) AdminHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Admin"))
}
