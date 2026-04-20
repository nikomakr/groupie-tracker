# Groupie Trackers

## Description

Groupie Trackers is a project that consists of receiving a given API and manipulating the data contained in it to create a site displaying the information. The API is made up of four parts:

- **Artists** — contains information about bands and artists: name(s), image, year they began their activity, date of their first album, and members
- **Locations** — contains their last and/or upcoming concert locations
- **Dates** — contains their last and/or upcoming concert dates
- **Relation** — links all the other parts together: artists, dates, and locations

## Objectives

- Build a user-friendly website that displays band/artist information through several data visualizations (blocks, cards, tables, lists, pages, graphics, etc.)
- Implement a client-server event: a feature of your choice that triggers an action requiring communication with the server to receive information (request-response)
- The website and server must not crash at any time
- All pages must work correctly and all errors must be handled

## Constraints

- The backend must be written in **Go**
- Only **standard Go packages** are allowed
- The server must handle errors correctly and never crash

## The API

Base URL: `https://groupietrackers.herokuapp.com/api`

| Endpoint | Returns |
|---|---|
| `/artists` | Array of artist objects |
| `/locations` | `{ "index": [ ... ] }` |
| `/dates` | `{ "index": [ ... ] }` |
| `/relation` | `{ "index": [ ... ] }` |

### Sample API responses

**`/artists[0]`**
```json
{
  "id": 1,
  "image": "https://groupietrackers.herokuapp.com/api/images/queen.jpeg",
  "name": "Queen",
  "members": ["Freddie Mercury", "Brian May", "John Daecon", "Roger Meddows-Taylor", "Mike Grose", "Barry Mitchell", "Doug Fogie"],
  "creationDate": 1970,
  "firstAlbum": "14-12-1973",
  "locations": "https://groupietrackers.herokuapp.com/api/locations/1",
  "concertDates": "https://groupietrackers.herokuapp.com/api/dates/1",
  "relations": "https://groupietrackers.herokuapp.com/api/relation/1"
}
```

**`/locations index[0]`**
```json
{
  "id": 1,
  "locations": ["north_carolina-usa", "georgia-usa", "los_angeles-usa", "saitama-japan"],
  "dates": "https://groupietrackers.herokuapp.com/api/dates/1"
}
```

**`/dates index[0]`**
```json
{
  "id": 1,
  "dates": ["*23-08-2019", "*22-08-2019", "*20-08-2019"]
}
```

**`/relation index[0]`**
```json
{
  "id": 1,
  "datesLocations": {
    "georgia-usa": ["22-08-2019"],
    "los_angeles-usa": ["20-08-2019"]
  }
}
```

## Project Structure

```
groupie-tracker/
├── cmd/
│   └── server/
│       └── main.go                          # Entry point — routing, server startup
├── internal/
│   ├── api/
│   │   ├── client.go                        # HTTP client, API fetch functions, Contains helper
│   │   └── client_test.go                   # Unit tests for the api package
│   ├── handlers/
│   │   └── handlers/
│   │       └── handlers.go                  # Exported HTTP handler functions
│   └── models/
│       └── models.go                        # Data model structs
├── templates/
│   ├── index.html                           # Homepage — artist card grid
│   └── artist.html                          # Artist detail page
└── go.mod
```

## Data Models

Defined in `internal/models/models.go`:

```go
type Artist struct {
    ID           int      `json:"id"`
    Image        string   `json:"image"`
    Name         string   `json:"name"`
    Members      []string `json:"members"`
    CreationDate int      `json:"creationDate"`
    FirstAlbum   string   `json:"firstAlbum"`
    Locations    string   `json:"locations"`
    ConcertDates string   `json:"concertDates"`
    Relations    string   `json:"relations"`
}

type Location struct {
    ID        int      `json:"id"`
    Locations []string `json:"locations"`
    Dates     string   `json:"dates"`
}

type LocationsIndex struct {
    Index []Location `json:"index"`
}

type Date struct {
    ID    int      `json:"id"`
    Dates []string `json:"dates"`
}

type DatesIndex struct {
    Index []Date `json:"index"`
}

type Relation struct {
    ID             int                 `json:"id"`
    DatesLocations map[string][]string `json:"datesLocations"`
}

type RelationsIndex struct {
    Index []Relation `json:"index"`
}
```

## API Client

Defined in `internal/api/client.go`. All requests go through a shared `http.Client` with a 15-second timeout.

### Individual endpoint functions

```go
api.GetArtists()   // returns ([]models.Artist, error)
api.GetLocations() // returns (models.LocationsIndex, error)
api.GetDates()     // returns (models.DatesIndex, error)
api.GetRelations() // returns (models.RelationsIndex, error)
```

Each function makes a GET request, reads the response body, unmarshals the JSON into the corresponding struct, and returns the data or an error.

### Concurrent batch fetch

```go
api.GetAllData() // returns ([]models.Artist, LocationsIndex, DatesIndex, RelationsIndex, error)
```

`GetAllData` fetches all four endpoints concurrently using goroutines and a `sync.WaitGroup`. If any request fails, the first error is returned and the remaining results are discarded. Used by the artist detail handler to reduce total latency.

### Search helper

```go
api.Contains(s, substr string) bool
```

Case-insensitive substring check. Used by the search handler to match artist names and member names against the query string.

## HTTP Handlers

Handler logic lives in two places:

- **`internal/handlers/handlers/handlers.go`** — exported handler functions (`IndexHandler`, `ArtistHandler`, `SearchHandler`) and their page data structs (`PageData`, `ArtistPageData`). `ArtistHandler` uses `api.GetAllData()` to fetch all four endpoints concurrently.
- **`cmd/server/main.go`** — wires routes to handlers and starts the server on `:8080`.

## Client-Server Event — Live Search

The search bar on the homepage is the implemented client-server event.

**Flow:**
1. User types in the search input (keyboard event)
2. Browser sends `GET /search?q=<query>` to the Go server
3. Server calls `api.GetArtists()` and filters results using `api.Contains()` (case-insensitive match on artist name and each member name)
4. Server responds with a JSON array of matching artists
5. Browser renders results as a dropdown — no page reload

**Why this qualifies as a client-server event:**
The client (browser) triggers an action (typing), which sends a request to the server, the server processes it and responds with data, and the client renders that data — a full request-response cycle.

## Pages

| Route | Description |
|---|---|
| `/` | Homepage — artist card grid |
| `/artist?id=N` | Artist detail — members, dates, locations, dates by location table |
| `/search?q=query` | JSON endpoint — returns matching artists |

## Error Handling

| Status | Trigger |
|---|---|
| `400` | Missing or invalid query parameter |
| `404` | Unknown route or artist ID not found |
| `405` | Wrong HTTP method |
| `500` | API fetch failure or template error |

## Progress

- [x] Project initialised — `go.mod` created
- [x] Data models defined — structs for Artist, Location, Date, Relation (`internal/models/models.go`)
- [x] API client — all four endpoints fetched and verified (`internal/api/client.go`)
- [x] Concurrent fetch — `GetAllData` fetches all four endpoints in parallel with goroutines
- [x] Search helper — `Contains` case-insensitive substring match in `internal/api`
- [x] HTTP server — routing, handlers, error handling (`cmd/server/main.go`)
- [x] Handlers refactored — exported `IndexHandler`, `ArtistHandler`, `SearchHandler` in `internal/handlers/handlers`
- [x] Homepage — artist card grid with dark theme (`templates/index.html`)
- [x] Artist detail page — members, dates, locations, relation table (`templates/artist.html`)
- [x] Client-server event — live search via `/search?q=`
- [x] Unit tests — `TestContains` and `TestContainsEmpty` in `internal/api/client_test.go`

## Usage

```bash
go run cmd/server/main.go
```

Visit `http://localhost:8080`

## Run Tests

```bash
go test ./internal/api/... -v
```

## Allowed Packages

Only Go standard library. No external dependencies. Verified via `go.mod`.

## Learning Outcomes

- Manipulation and storage of data
- JSON files and format
- HTML
- Event creation and visualization
- Client-server communication