# Groupie Trackers — Audit Checklist

## Functional

### Has the requirement for allowed packages been respected?
Yes. Only Go standard library packages are used:
`net/http`,
`encoding/json`,
`html/template`,
`fmt`,
`log`,
`io`,
`net/http`,
`strconv`,
`strings`,
`time`.

Verified in `go.mod` — zero external dependencies.

### Is the data from `artists` being used?
Yes. Fetched via `api.GetArtists()`. Displayed on the homepage as cards (name, image, creation date) and on the artist detail page (name, image, creation date, first album, members).

Test it by running the programme:
go run cmd/server/main.go

Then run the following command:
curl -s http://localhost:8080 | grep -o '<strong>[^<]*</strong>' | head -5

This hits your live server and pulls the artist names out of the rendered HTML. If you see names like <strong>Queen</strong> printed back, it proves the artists data is being fetched and rendered.

Expected outcome:
<strong>Queen</strong>
<strong>SOJA</strong>
<strong>Pink Floyd</strong>
<strong>Scorpions</strong>
<strong>XXXTentacion</strong>


### Is the data from `locations` being used?
Yes. Fetched via `api.GetLocations()`. Displayed on the artist detail page as a list of concert locations.

Run following command:
curl -s "http://localhost:8080/artist?id=1" | grep -o '<li>[^<]*</li>' | head -10

This hits the artist detail page for Queen and pulls out the list items — which include the locations. If you see location names printed back, it proves locations data is being fetched and rendered.

Expected outcome:
<li>Freddie Mercury</li>
<li>Brian May</li>
<li>John Daecon</li>
<li>Roger Meddows-Taylor</li>
<li>Mike Grose</li>
<li>Barry Mitchell</li>
<li>Doug Fogie</li>
<li>north_carolina-usa</li>
<li>georgia-usa</li>
<li>los_angeles-usa</li>

### Is the data from `dates` being used?
Yes. Fetched via `api.GetDates()`. Displayed on the artist detail page as a list of concert dates.

Test it by running the following command:
curl -s "http://localhost:8080/artist?id=1" | grep -o '<li>[^<]*</li>' | grep -E '[0-9]{2}-[0-9]{2}-[0-9]{4}'

This filters the list items to only show ones that match a date pattern like 23-08-2019.

Expected outcome:
<li>*07-02-2020</li>
<li>*10-02-2020</li>
<li>dunedin-new_zealand: 10-02-2020 </li>
<li>georgia-usa: 22-08-2019 </li>
<li>los_angeles-usa: 20-08-2019 </li>
<li>nagoya-japan: 30-01-2019 </li>
<li>north_carolina-usa: 23-08-2019 </li>
<li>osaka-japan: 28-01-2020 </li>
<li>penrose-new_zealand: 07-02-2020 </li>
<li>saitama-japan: 26-01-2020 </li>

### Is data from `relations` being used?
Yes. Fetched via `api.GetRelations()`. Displayed on the artist detail page as a table mapping each location to its concert dates.

Test it by running the following command:
curl -s "http://localhost:8080/artist?id=1" | grep -i "georgia\|london\|japan"

Expected outcome:
<li>georgia-usa</li>
<li>saitama-japan</li>
<li>osaka-japan</li>
<li>nagoya-japan</li>
<li>georgia-usa: 22-08-2019 </li>
<li>nagoya-japan: 30-01-2019 </li>
<li>osaka-japan: 28-01-2020 </li>
<li>saitama-japan: 26-01-2020 </li>

### Members for Queen — does it show the right members?
Navigate to `http://localhost:8080/artist?id=1`. The Members section displays:
```
Freddie Mercury
Brian May
John Daecon
Roger Meddows-Taylor
Mike Grose
Barry Mitchell
Doug Fogie
```


### First album for Gorillaz — does it show the right date?
Navigate to `http://localhost:8080/artist?id=2`. The First Album field displays:
```
26-03-2001
```

Alternatively, you can run the following command:
curl -s "https://groupietrackers.herokuapp.com/api/artists" | python3 -c "import json,sys; artists=json.load(sys.stdin); [print(a['firstAlbum']) for a in artists if a['name']=='Gorillaz']"

### Locations for Travis Scott — does it show the right locations?
Navigate to Travis Scott's detail page. The Locations section displays the raw API values:
```
santiago-chile
sao_paulo-brasil
los_angeles-usa
houston-usa
atlanta-usa
new_orleans-usa
philadelphia-usa
london-uk
frauenfeld-switzerland
turku-finland
```


### Members for Foo Fighters — does it show the right members?
Navigate to Foo Fighters' detail page. The Members section displays:
```
Dave Grohl
Nate Mendel
Taylor Hawkins
Chris Shiflett
Pat Smear
Rami Jaffee
```


### Does a client-server event exist and respond as expected?
Yes. The search bar on the homepage is the event:
1. User types in the search input
2. Browser sends `GET /search?q=<query>` to the Go server
3. Server searches artists by name and members (case-insensitive)
4. Server responds with JSON
5. Browser renders a dropdown of results — no page reload

Test it: type "queen" or "freddie" in the search bar on the homepage.


### Did the server behave as expected and not crash?
Yes. Every handler checks for errors and returns the appropriate HTTP response. The server never panics.

Test it:
```bash
curl -s http://localhost:8080/artist?id=99999
curl -s http://localhost:8080/notapage
curl -X POST http://localhost:8080/
```

Expected responses:
- `404 - Artist Not Found`
- `404 - Page Not Found`
- `405 - Method Not Allowed`

The server keeps running after all three — it never crashes.


### Does the server use the right HTTP method?
Yes. All routes use `GET`. Any other method returns `405 - Method Not Allowed`. Test with:
```bash
curl -X POST http://localhost:8080/
```


### Did the site run without crashing?
Yes. All error paths return HTTP responses instead of crashing. The server runs until manually stopped.

Test it by hitting bad routes back to back while the server is running:
```bash
curl -s http://localhost:8080/notapage
curl -s http://localhost:8080/artist
curl -s http://localhost:8080/artist?id=abc
curl -s http://localhost:8080/artist?id=99999
curl -X DELETE http://localhost:8080/
```

Then visit `http://localhost:8080` in the browser — the homepage still loads normally. The server absorbed every bad request and kept running.


### Are all pages working — no 404?
Yes. Defined routes:
- `/` — Homepage
- `/artist?id=N` — Artist detail
- `/search?q=query` — Live search JSON

Any undefined path returns a proper `404 - Page Not Found` response.


### Does the project handle HTTP 500 — Internal Server Errors?
Yes. Every handler wraps API calls and template execution in error checks. On failure it returns `500 - Internal Server Error` and logs the reason.

Two failure points are handled in every handler:
1. API fetch failure — network error, timeout, bad response
2. Template parse or execute failure — missing or malformed template file

In both cases the server responds with 500 and logs the reason without crashing.


### Is communication between server and client well established?
Yes. Two patterns:
1. Server-side rendering — Go templates render full HTML on request
2. AJAX — `/search` returns `Content-Type: application/json` consumed by `fetch()` in the browser

Run this to prove both patterns:

Pattern 1 — server-side rendering:
bashcurl -s http://localhost:8080 | grep "<title>"
Should return <title>Groupie Trackers</title> — full HTML rendered by the server.
Pattern 2 — AJAX JSON response:
bashcurl -s "http://localhost:8080/search?q=queen"
Should return a JSON array — consumed by fetch() in the browser without a page reload.
Both prove the server is communicating correctly with the client in two different ways.


### Does the server present all needed handlers and patterns?
Yes. Three handlers registered in `main()`:
- `IndexHandler` — serves the homepage
- `ArtistHandler` — serves the artist detail page
- `SearchHandler` — handles the client-server event


## General

### Does the event system run asynchronously?
Yes. Explicitly implemented with goroutines and `sync.WaitGroup`.

When the artist detail page is requested, all four API endpoints are fetched concurrently:
```go
wg.Add(4)

go func() { defer wg.Done(); errs[0] = fetchJSON(baseURL+"/artists", &artists) }()
go func() { defer wg.Done(); errs[1] = fetchJSON(baseURL+"/locations", &locations) }()
go func() { defer wg.Done(); errs[2] = fetchJSON(baseURL+"/dates", &dates) }()
go func() { defer wg.Done(); errs[3] = fetchJSON(baseURL+"/relation", &relations) }()

wg.Wait()
```

All four goroutines run in parallel. `wg.Wait()` blocks until all four finish. This means the page loads in the time of the slowest single request, not the sum of all four.

On the client side, the search event uses `fetch()` which is non-blocking — the browser sends the request and continues running without freezing the page.

### Is the site hosted/deployed? Can you access the website through a DNS (Domain Name System)?

Yes. The site is deployed on Render and publicly accessible at:

https://groupie-tracker-h5ic.onrender.com/

Deployed from the `main` branch of https://github.com/nikomakr/groupie-tracker. Render builds and runs the server automatically on every push.

For local development: `go run cmd/server/main.go` → http://localhost:8080

## Basic

### Does the project run quickly and effectively?
Yes. Two things prove it:

**1. Concurrent fetching on the artist detail page**

All four API endpoints are fetched in parallel using goroutines. Instead of waiting for each one sequentially, they all run at the same time. The page loads in the time of the slowest single request, not the sum of all four.

**2. 15-second timeout on the HTTP client**
```go
var httpClient = &http.Client{Timeout: 15 * time.Second}
```

If the upstream API hangs, the server gives up after 15 seconds and returns a 500 instead of waiting forever.

To prove the speed, time the artist detail page:
```bash
time curl -s "http://localhost:8080/artist?id=1" > /dev/null
```

You should see a response in under 1 second.


### Does the code obey good practices?
Yes. Four areas demonstrate this:

**1. Packages separated by responsibility**
- `internal/models` — data structures only, no logic
- `internal/api` — all external API communication
- `internal/handlers` — all HTTP request handling
- `cmd/server/main.go` — server setup and routing only

**2. Errors are always handled**
No error is silently ignored. Every function returns an error and every caller checks it:
```go
artists, err := api.GetArtists()
if err != nil {
    http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
    log.Println("Error fetching artists:", err)
    return
}
```

**3. Exported vs unexported identifiers used correctly**
- `Contains` is exported — used by both `handlers` and `client_test.go`
- `fetchJSON` is unexported — internal helper, not needed outside `api` package

**4. HTTP timeouts configured**
Both the HTTP client (15s) and all error paths are handled so the server never hangs or crashes.

Verify with:
```bash
go build ./...
go test ./internal/api/... -v
```
Zero errors and passing tests confirm the code is clean.


### Is there a test file?
Yes. `internal/api/client_test.go` contains unit tests for the `Contains` function. Run with:
```bash
go test ./internal/api/... -v
```

---

## Social

### Did you learn anything from this project?
Yes. This project covered:

- **Consuming a REST API in Go** — making HTTP GET requests, reading response bodies, handling timeouts
- **JSON** — unmarshaling API responses into typed Go structs using `json` tags
- **HTML templating** — rendering dynamic pages with `html/template`, passing data from Go into HTML
- **Client-server events** — triggering a server request from the browser using `fetch()` and handling the JSON response without a page reload
- **Concurrency** — fetching multiple API endpoints in parallel using goroutines and `sync.WaitGroup`
- **Error handling** — returning correct HTTP status codes (400, 404, 405, 500) and never crashing the server
- **Package structure** — separating concerns into `models`, `api`, `handlers` and keeping `main.go` clean
- **Go tooling** — `go mod init`, `go build`, `go run`, `go test`

### Can it be open-sourced / used for other purposes?
Yes. The structure is generic enough to be reused for any Go web application that consumes an external API. Specifically:

- `internal/models` — swap the structs for any API schema
- `internal/api` — swap the endpoints for any REST API
- `internal/handlers` — the error handling and routing pattern applies to any project
- `templates/` — the dark theme and card grid layout can be reused for any data display

The concurrent fetching pattern with goroutines and `sync.WaitGroup` is directly reusable in any Go project that needs to fetch multiple endpoints at the same time.

### Would you nominate this as an example project?
Yes. It demonstrates several key principles that make it a good reference:

- **Zero dependencies** — built entirely with the Go standard library, no external packages
- **Clean structure** — concerns are separated into `models`, `api`, `handlers` with `main.go` only handling routing
- **Explicit concurrency** — goroutines and `sync.WaitGroup` are used intentionally, not just relied on implicitly
- **Full error coverage** — every possible failure point returns the correct HTTP status code
- **Real client-server event** — not just a page link, but a live `fetch()` request to a Go JSON endpoint that responds without a page reload
- **Tested** — unit tests exist and pass

Any of these points alone is good practice. All of them together in a single project built from scratch makes it worth nominating.