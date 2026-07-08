// Vercel entry point
package api

import (
	"net/http"

	"pitch.ideas/internal/app"
)


var router = app.NewRouter()

func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w,r)
}