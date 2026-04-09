package main

import (
	"html/template"
	"log"
	"net/http"

	"groupie-tracker/internal/api"
	"groupie-tracker/internal/models"
)

type PageData struct {
	Artists []models.Artist
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)

	log.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 - Page Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	artists, err := api.GetArtists()
	if err != nil {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		log.Println("Error fetching artists:", err)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		log.Println("Error parsing template:", err)
		return
	}

	data := PageData{Artists: artists}
	tmpl.Execute(w, data)
}