package web

import (
	"html/template"
	"net/http"
)

func render(w http.ResponseWriter, tmpl string, data any) {
	t, err := template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/web/"+tmpl,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
