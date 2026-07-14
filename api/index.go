// Vercel entry point
package handler

import (
	"net/http"
	"fmt"
	"os"
	"sync"

	"pitch.ideas/pkg/server"
)

var (
	router http.Handler
	once sync.Once
)

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		fmt.Println("Init called")
		router = server.NewRouter()
	})

    fmt.Println("Handler called")
	router.ServeHTTP(w, r)
}