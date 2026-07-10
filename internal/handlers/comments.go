package handlers

import (
	"net/http"
	"encoding/json"

	"strconv"
	"github.com/go-chi/chi/v5"

	"pitch.ideas/internal/auth"
	"pitch.ideas/internal/database"
)


type CommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

func AddComment(w http.ResponseWriter, r *http.Request) {
	idea_id, err := strconv.Atoi(chi.URLParam(r, "idea_id"))
	if err != nil {
		http.Error(w, "invalid idea id", http.StatusBadRequest)
		return
	}
	ideaID := uint(idea_id)

		var req CommentRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

	user := auth.GetUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = database.CreateComment(ideaID, user.ID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func EditComment(w http.ResponseWriter, r *http.Request) {
	comment_id, err := strconv.Atoi(chi.URLParam(r, "comment_id"))
	if err != nil {
		http.Error(w, "invalid idea id", http.StatusBadRequest)
		return
	}
	commentID := uint(comment_id)

		var req CommentRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

	user := auth.GetUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = database.EditComment(commentID, user.ID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	comment_id, err := strconv.Atoi(chi.URLParam(r, "comment_id"))
	if err != nil {
		http.Error(w, "invalid idea id", http.StatusBadRequest)
		return
	}
	commentID := uint(comment_id)

	user := auth.GetUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = database.DeleteComment(commentID, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
