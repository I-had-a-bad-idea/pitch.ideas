package auth

import (
	"net/http"
	"context"
	
	"pitch.ideas/internal/database"
	"pitch.ideas/internal/database/models"
)


const userContextKey string = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user := database.GetUserBySession(cookie.Value)
		if user == nil {
			http.Error(w, "Expired", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


func GetUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok {
		return nil
	}

	return user
}