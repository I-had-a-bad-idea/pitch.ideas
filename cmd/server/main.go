package main

import (
	"fmt"
	"log"

	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"pitch.ideas/internal/handlers"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", handlers.Home)


	fmt.Println("Server running on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", r))

	fmt.Println("hello")
}  