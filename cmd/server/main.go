package main

import (
	"log"
	"net/http"
	"os"

	"groupie-tracker/internal/handlers/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.IndexHandler)
	mux.HandleFunc("/artist", handlers.ArtistHandler)
	mux.HandleFunc("/search", handlers.SearchHandler)

	log.Printf("Server running on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}