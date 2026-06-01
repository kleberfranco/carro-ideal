package web

import "net/http"

func (h *Handler) RecommendHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "recommend.html", map[string]any{
		"Title": "Recomendar Carro",
	})
}
