package main

import (
	"fmt"
	"log" // for logging errors

	"groupie-tracker/internal/api"
)

func main() {
	artists, err := api.GetArtists()
	if err != nil {
		log.Fatal(err) // log the error and exit if fetching artists fails
	}
	fmt.Println("First artist:", artists[0].Name)
	fmt.Println("Members:", artists[0].Members)

	locations, err := api.GetLocations()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("First location:", locations.Index[0].Locations)

	dates, err := api.GetDates()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("First dates:", dates.Index[0].Dates)

	relations, err := api.GetRelations()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("First relation:", relations.Index[0].DatesLocations)
}