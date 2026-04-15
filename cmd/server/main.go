package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"groupie-tracker/internal/api"
	"groupie-tracker/internal/models"
)

type PageData struct {
	Artists []models.Artist
}

type ArtistPageData struct {
	Artist   models.Artist
	Location models.Location
	Date     models.Date
	Relation models.Relation
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/artist", artistHandler)
	mux.HandleFunc("/search", searchHandler)

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

func artistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "400 - Bad Request", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "400 - Bad Request", http.StatusBadRequest)
		return
	}

	artists, err := api.GetArtists()
	if err != nil {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		log.Println("Error fetching artists:", err)
		return
	}

	locations, err := api.GetLocations()
	if err != nil {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		log.Println("Error fetching locations:", err)
		return
	}

	dates, err := api.GetDates()
	if err != nil {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		log.Println("Error fetching dates:", err)
		return
	}

	relations, err := api.GetRelations()
	if err != nil {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		log.Println("Error fetching relations:", err)
		return
	}

	var artist models.Artist
	for _, a := range artists {
		if a.ID == id {
			artist = a
			break
		}
	}
	if artist.ID == 0 {
		http.Error(w, "404 - Artist Not Found", http.StatusNotFound)
		return
	}

	var location models.Location
	for _, l := range locations.Index {
		if l.ID == id {
			location = l
			break
		}
	}

	var date models.Date
	for _, d := range dates.Index {
		if d.ID == id {
			date = d
			break
		}
	}

	var relation models.Relation
	for _, rel := range relations.Index {
		if rel.ID == id {
			relation = rel
			break
		}
	}

	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		log.Println("Error parsing template:", err)
		return
	}

	data := ArtistPageData{
		Artist:   artist,
		Location: location,
		Date:     date,
		Relation: relation,
	}
	tmpl.Execute(w, data)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "400 - Bad Request", http.StatusBadRequest)
		return
	}

	artists, err := api.GetArtists()
	if err != nil {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		log.Println("Error fetching artists:", err)
		return
	}

	var results []models.Artist
	for _, a := range artists {
		if contains(a.Name, query) {
			results = append(results, a)
			continue
		}
		for _, m := range a.Members {
			if contains(m, query) {
				results = append(results, a)
				break
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}