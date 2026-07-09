package app

import (

	"net/http"
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"pitch.ideas/internal/handlers"
	"pitch.ideas/internal/views"
	"pitch.ideas/internal/database"
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

func NewRouter() http.Handler {
	renderer := views.New()
	database.InitDB()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/", handlers.Home(renderer))

	r.Route("/auth", func(r chi.Router) {
		r.Get("/login", handlers.LoginPage(renderer))
		r.Post("/login", handlers.Login)

		r.Get("/register", handlers.RegisterPage(renderer))
		r.Post("/register", handlers.Register)

		r.Get("/logout", handlers.LogoutPage(renderer))
		r.Post("/logout", handlers.Logout)

		r.Get("/status", handlers.AuthStatus)
	})

	r.Route("/pitches", func(r chi.Router) {
		r.Get("/", handlers.ListPitches)

		r.With(AuthMiddleware).Get("/create", handlers.CreatePitchPage(renderer))
		r.With(AuthMiddleware).Put("/create", handlers.CreatePitch)

		r.Get("/{idea_id}", handlers.GetPitchPage(renderer))
		r.With(AuthMiddleware).Post("/{idea_id}/upvote", handlers.UpvotePitch)

		r.With(AuthMiddleware).Post("/{idea_id}/edit", handlers.EditPitch)
		r.With(AuthMiddleware).Delete("/{idea_id}/delete", handlers.DeletePitch)

		r.Route("/{idea_id}/comments", func(r chi.Router) {
			r.With(AuthMiddleware).Post("/add", handlers.AddComment)

			r.With(AuthMiddleware).Post("/{comment_id}/edit", handlers.EditComment)
			r.With(AuthMiddleware).Delete("/{comment_id}/delete", handlers.DeleteComment)
		})
	})

	return  r
}  