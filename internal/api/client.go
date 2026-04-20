package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"groupie-tracker/internal/models"
)

const baseURL = "https://groupietrackers.herokuapp.com/api"

var httpClient = &http.Client{Timeout: 15 * time.Second}

func fetchJSON(url string, dest interface{}) error {
	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("fetching %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading body from %s: %w", url, err)
	}

	if err := json.Unmarshal(body, dest); err != nil {
		return fmt.Errorf("decoding %s: %w", url, err)
	}
	return nil
}

func GetAllData() ([]models.Artist, models.LocationsIndex, models.DatesIndex, models.RelationsIndex, error) {
	var (
		artists   []models.Artist
		locations models.LocationsIndex
		dates     models.DatesIndex
		relations models.RelationsIndex
		errs      = make([]error, 4)
		wg        sync.WaitGroup
	)

	wg.Add(4)

	go func() {
		defer wg.Done()
		errs[0] = fetchJSON(baseURL+"/artists", &artists)
	}()

	go func() {
		defer wg.Done()
		errs[1] = fetchJSON(baseURL+"/locations", &locations)
	}()

	go func() {
		defer wg.Done()
		errs[2] = fetchJSON(baseURL+"/dates", &dates)
	}()

	go func() {
		defer wg.Done()
		errs[3] = fetchJSON(baseURL+"/relation", &relations)
	}()

	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return nil, models.LocationsIndex{}, models.DatesIndex{}, models.RelationsIndex{}, err
		}
	}

	return artists, locations, dates, relations, nil
}

func GetArtists() ([]models.Artist, error) {
	var artists []models.Artist
	if err := fetchJSON(baseURL+"/artists", &artists); err != nil {
		return nil, err
	}
	return artists, nil
}

func GetLocations() (models.LocationsIndex, error) {
	var locations models.LocationsIndex
	if err := fetchJSON(baseURL+"/locations", &locations); err != nil {
		return models.LocationsIndex{}, err
	}
	return locations, nil
}

func GetDates() (models.DatesIndex, error) {
	var dates models.DatesIndex
	if err := fetchJSON(baseURL+"/dates", &dates); err != nil {
		return models.DatesIndex{}, err
	}
	return dates, nil
}

func GetRelations() (models.RelationsIndex, error) {
	var relations models.RelationsIndex
	if err := fetchJSON(baseURL+"/relation", &relations); err != nil {
		return models.RelationsIndex{}, err
	}
	return relations, nil
}

func Contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}