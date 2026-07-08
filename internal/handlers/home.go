package handlers

import (
	"net/http"

	"pitch.ideas/internal/handlers/views"
)


func Home(renderer *views.Renderer) http.HandlerFunc {
	return  func(w http.ResponseWriter, r *http.Request) {
		renderer.Render(w, "index.html", "")
	}
}