package handlers

import (
	"net/http"

	"pitch.ideas/internal/handlers/views"
)

func RegisterPage(renderer *views.Renderer) http.HandlerFunc {
	return  func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Register page"))
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
}