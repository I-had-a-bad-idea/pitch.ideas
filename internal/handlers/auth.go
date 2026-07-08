package handlers

import (
	"net/http"

	"pitch.ideas/internal/views"
)

func LoginPage(renderer *views.Renderer) http.HandlerFunc {
	return  func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login page"))
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
}


func RegisterPage(renderer *views.Renderer) http.HandlerFunc {
	return  func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Register page"))
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
}