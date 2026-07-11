// Vercel entry point
package handler

import (
	"net/http"
	"fmt"

	"pitch.ideas/pkg/server"
)

var router http.Handler

func init() {
    fmt.Println("Init called")
	router = server.NewRouter()
}

func Handler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Handler called")
	router.ServeHTTP(w, r)
}