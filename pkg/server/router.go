package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"pitch.ideas/internal/auth"
	"pitch.ideas/internal/database"
	"pitch.ideas/internal/handlers"
	"pitch.ideas/internal/views"
)



func NewRouter() http.Handler {
	renderer := views.New()
	if err := database.InitDB(); err != nil {
		log.Fatalf("DB init failed: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	fs := http.FileServer(http.Dir("public"))
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

		r.With(auth.AuthMiddleware).Get("/create", handlers.CreatePitchPage(renderer))
		r.With(auth.AuthMiddleware).Put("/create", handlers.CreatePitch)

		r.Get("/{idea_id}", handlers.GetPitchPage(renderer))
		r.With(auth.AuthMiddleware).Post("/{idea_id}/upvote", handlers.UpvotePitch)

		r.With(auth.AuthMiddleware).Post("/{idea_id}/edit", handlers.EditPitch)
		r.With(auth.AuthMiddleware).Delete("/{idea_id}/delete", handlers.DeletePitch)

		r.Route("/{idea_id}/comments", func(r chi.Router) {
			r.With(auth.AuthMiddleware).Post("/add", handlers.AddComment)

			r.With(auth.AuthMiddleware).Post("/{comment_id}/edit", handlers.EditComment)
			r.With(auth.AuthMiddleware).Delete("/{comment_id}/delete", handlers.DeleteComment)
		})
	})

	return  r
}  