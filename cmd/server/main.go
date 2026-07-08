package main

import (
	"fmt"
	"log"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"pitch.ideas/internal/handlers"
	"pitch.ideas/internal/views"
)

func main() {
	renderer := views.New()

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
	})

	r.Route("/pitches", func(r chi.Router) {
		r.Get("/", handlers.ListPitches)

		r.Get("/create", handlers.CreatePitchPage(renderer))
		r.Put("/create", handlers.CreatePitch)

		r.Get("/{id}", handlers.GetPitch)
		r.Post("/{id}/edit", handlers.EditPitch)
		r.Delete("/{id}/delete", handlers.DeletePitch)
	})

	fmt.Println("Server running on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", r))

	fmt.Println("hello")
}  