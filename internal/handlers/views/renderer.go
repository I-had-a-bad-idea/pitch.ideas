package views

import (
	"html/template"
	"net/http"
)

type Renderer struct {
	templates *template.Template
}

func New() *Renderer {
	return &Renderer{
		templates: template.Must(
			template.ParseGlob("templates/*.html"),
		),
	}
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data any) {
	err := r.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}