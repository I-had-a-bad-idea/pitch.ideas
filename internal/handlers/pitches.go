package handlers

import (
	"net/http"
	"encoding/json"
	
	"pitch.ideas/internal/views"
	"pitch.ideas/internal/database"
)

func ListPitches(w http.ResponseWriter, r *http.Request) {
	pitches, err := database.GetAllIdeasAsDicts(20)
	
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]any{
		"pitches": pitches,
	}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
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