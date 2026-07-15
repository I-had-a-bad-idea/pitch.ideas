package main
// Local server
import (
	"log"
	"net/http"

	"pitch.ideas/pkg/server"
)

func main() {
	router := server.NewRouter()

	log.Println("Server running on http://localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", router))
}  
