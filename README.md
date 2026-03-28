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

## Usage

```bash
go run .
```

Visit `http://localhost:8080`

## Learning Outcomes

- Manipulation and storage of data
- JSON files and format
- HTML
- Event creation and visualization
- Client-server communication