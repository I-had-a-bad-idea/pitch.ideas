package views

import (
	"html/template"
	"net/http"
	"embed"
)

type Renderer struct {
	templates *template.Template
}

//go:embed templates/*.html
var templates embed.FS

func New() *Renderer {
	return &Renderer{
		templates: template.Must(
			template.ParseFS(templates, "templates/*.html"),
		),
	}
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data any) {
	err := r.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}