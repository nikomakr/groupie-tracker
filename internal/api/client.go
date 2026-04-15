package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"strings"

	"groupie-tracker/internal/models"
)

const baseURL = "https://groupietrackers.herokuapp.com/api"

var httpClient = &http.Client{Timeout: 15 * time.Second} // just for safety, to avoid hanging requests

func GetArtists() ([]models.Artist, error) {
	resp, err := httpClient.Get(baseURL + "/artists")
	if err != nil {
		return nil, fmt.Errorf("fetching artists: %w", err)
	}
	defer resp.Body.Close() // ensure the response body is closed after we're done with it

	body, err := io.ReadAll(resp.Body) // read the entire response body
	if err != nil {
		return nil, fmt.Errorf("reading artists body: %w", err)
	}

	var artists []models.Artist // unmarshal the JSON response into a slice of Artist structs
	if err := json.Unmarshal(body, &artists); err != nil { // handle JSON decoding errors
		return nil, fmt.Errorf("decoding artists: %w", err)
	}
	return artists, nil
}

func GetLocations() (models.LocationsIndex, error) { // fetch the locations index from the API and return it as a LocationsIndex struct
	resp, err := httpClient.Get(baseURL + "/locations")
	if err != nil {
		return models.LocationsIndex{}, fmt.Errorf("fetching locations: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.LocationsIndex{}, fmt.Errorf("reading locations body: %w", err)
	}

	var locations models.LocationsIndex
	if err := json.Unmarshal(body, &locations); err != nil {
		return models.LocationsIndex{}, fmt.Errorf("decoding locations: %w", err)
	}
	return locations, nil
}

func GetDates() (models.DatesIndex, error) {
	resp, err := httpClient.Get(baseURL + "/dates")
	if err != nil {
		return models.DatesIndex{}, fmt.Errorf("fetching dates: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.DatesIndex{}, fmt.Errorf("reading dates body: %w", err)
	}

	var dates models.DatesIndex
	if err := json.Unmarshal(body, &dates); err != nil {
		return models.DatesIndex{}, fmt.Errorf("decoding dates: %w", err)
	}
	return dates, nil
}

func GetRelations() (models.RelationsIndex, error) {
	resp, err := httpClient.Get(baseURL + "/relation")
	if err != nil {
		return models.RelationsIndex{}, fmt.Errorf("fetching relations: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.RelationsIndex{}, fmt.Errorf("reading relations body: %w", err)
	}

	var relations models.RelationsIndex
	if err := json.Unmarshal(body, &relations); err != nil {
		return models.RelationsIndex{}, fmt.Errorf("decoding relations: %w", err)
	}
	return relations, nil
}

func Contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}