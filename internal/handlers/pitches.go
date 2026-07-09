package handlers

import (
	"net/http"
	"encoding/json"

	"strconv"
	"github.com/go-chi/chi/v5"

	"pitch.ideas/internal/views"
	"pitch.ideas/internal/database"
	"pitch.ideas/internal/auth"
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

type PitchRequest struct {
    Title       string `json:"title" validate:"required,min=1,max=100"`
    Topic       string `json:"topic" validate:"required,min=1,max=50"`
	Description string `json:"description" validate:"required,min=1,max=5000"`
}

func CreatePitch(w http.ResponseWriter, r *http.Request) {
	var req PitchRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user := auth.GetUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}


	database.CreateIdea(req.Title, req.Topic, req.Description, user.ID)
}

func GetPitchPage(renderer *views.Renderer) http.HandlerFunc {
	return  func(w http.ResponseWriter, r *http.Request) {
		ideaID, err := strconv.Atoi(chi.URLParam(r, "idea_id"))
		if err != nil {
			http.Error(w, "invalid idea id", http.StatusBadRequest)
			return
		}

		idea, err := database.GetIdeaDict(uint(ideaID))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if idea == nil {
			http.Error(w, "Pitch not found", http.StatusNotFound)
			return 
		}

		comments, err := database.GetCommentsDict(uint(ideaID), 50)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		renderer.Render(w, "pitch.html", map[string]any{
			"idea": idea,
			"comments": comments,
		})
	}
}

func UpvotePitch(w http.ResponseWriter, r *http.Request) {
}

func EditPitch(w http.ResponseWriter, r *http.Request) {
}

func DeletePitch(w http.ResponseWriter, r *http.Request) {
}