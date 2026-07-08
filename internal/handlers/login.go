package handlers

import (
	"net/http"

	"pitch.ideas/internal/handlers/views"
)

func LoginPage(renderer *views.Renderer) http.HandlerFunc {
	return  func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login page"))
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
}