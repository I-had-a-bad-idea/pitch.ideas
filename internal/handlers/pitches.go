package handlers

import (
	"net/http"

	"pitch.ideas/internal/views"
)

func ListPitches(w http.ResponseWriter, r *http.Request) {
}

func CreatePitchPage(renderer *views.Renderer) http.HandlerFunc {
	return  func(w http.ResponseWriter, r *http.Request) {
		renderer.Render(w, "create-pitch.html", "")
	}
}

func CreatePitch(w http.ResponseWriter, r *http.Request) {
}

func GetPitchPage(renderer *views.Renderer) http.HandlerFunc {
	return  func(w http.ResponseWriter, r *http.Request) {
		renderer.Render(w, "pitch.html", "")
	}
}

func UpvotePitch(w http.ResponseWriter, r *http.Request) {
}

func EditPitch(w http.ResponseWriter, r *http.Request) {
}

func DeletePitch(w http.ResponseWriter, r *http.Request) {
}