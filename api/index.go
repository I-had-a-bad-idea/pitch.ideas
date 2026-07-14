// Vercel entry point
package handler

import (
	"net/http"
	"fmt"
	"sync"
	"log"
	"os"

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

	wd, _ := os.Getwd()
	log.Println("Working dir:", wd)

	entries, err := os.ReadDir("public")
	if err != nil {
		log.Println(err)
	} else {
		for _, e := range entries {
			log.Println(e.Name())
		}
	}

    fmt.Println("Handler called")
	router.ServeHTTP(w, r)
}