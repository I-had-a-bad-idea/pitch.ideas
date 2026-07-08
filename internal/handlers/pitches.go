package handlers

import (
	"net/http"

	"pitch.ideas/internal/views"
)

func ListPitches(w http.ResponseWriter, r *http.Request) {
}

func CreatePitchPage(renderer *views.Renderer) http.HandlerFunc {
	return  func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login page"))
	}
}

func CreatePitch(w http.ResponseWriter, r *http.Request) {
}

func GetPitch(w http.ResponseWriter, r *http.Request) {
}

func UpvotePitch(w http.ResponseWriter, r *http.Request) {
}

func EditPitch(w http.ResponseWriter, r *http.Request) {
}

func DeletePitch(w http.ResponseWriter, r *http.Request) {
}