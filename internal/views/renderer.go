package views

import (
	"html/template"
	"net/http"
	"path/filepath"
	"os"
)

type Renderer struct {
	templates *template.Template
}

func New() *Renderer {
	root, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return &Renderer{
		templates: template.Must(
			template.ParseGlob(
				filepath.Join(root, "templates", "*.html"),
			),
		),
	}
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data any) {
	err := r.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}