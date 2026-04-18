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
05-06-2002
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


### Does the server present all needed handlers and patterns?
Yes. Three handlers registered in `main()`:
- `IndexHandler` — serves the homepage
- `ArtistHandler` — serves the artist detail page
- `SearchHandler` — handles the client-server event

---

## General

### Does the event system run asynchronously?
The search event is triggered asynchronously from the browser via `fetch()` — it does not block the page. On the server side, each HTTP request is handled in its own goroutine by Go's `net/http` package automatically.

---

## Basic

### Does the project run quickly and effectively?
Yes. The API client uses a 15-second timeout. Each handler fetches only what it needs. No unnecessary data requests.

---

### Does the code obey good practices?
Yes:
- Packages separated by responsibility: `models`, `api`, `handlers`
- Errors are always handled and never silently ignored
- Exported vs unexported identifiers used correctly
- `main.go` only contains server setup and routing

---

### Is there a test file?
Yes. `internal/api/client_test.go` contains unit tests for the `Contains` function. Run with:
```bash
go test ./internal/api/... -v
```

---

## Social

### Did you learn anything from this project?
This project covers: consuming a REST API in Go, working with JSON, rendering HTML with `html/template`, creating client-server events using `fetch()` and a Go JSON endpoint, and structuring a Go web application cleanly.

### Can it be open-sourced / used for other purposes?
Yes. The structure — API client, handlers, models, templates — is a clean starting point for any Go web application that consumes and displays external API data.

### Would you nominate this as an example project?
Yes. It uses only the standard library, handles errors correctly, separates concerns cleanly, and implements a working client-server event.