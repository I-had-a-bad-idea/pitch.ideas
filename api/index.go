// Vercel entry point
package handler

import (
	"net/http"

	"pitch.ideas/pkg/server"
)


var router = server.NewRouter()

func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w,r)
}