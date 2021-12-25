package main

import (
	"net/http"
	"time"

	"go-movies/models"
)

func (app *application) getOneMovie(w http.ResponseWriter, r *http.Request) {
	id, err := app.getIdFromUrlParams(w, r)
	if err != nil {
		app.clientErrorJSON(w, err)
		return
	}
	movie := models.Movie{
		ID:          id,
		Title:       "some movie",
		Description: "some description",
		Year:        2021,
		ReleaseDate: time.Date(2021, 01, 01, 01, 0, 0, 0, time.Local),
		Runtime:     125,
		Rating:      5,
		MPAARating:  "PG-13",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err = app.writeJSON(w, http.StatusOK, movie, "movie"); err != nil {
		app.serverErrorJSON(w, err)
		return
	}
}

func (app *application) getAllMovies(w http.ResponseWriter, r *http.Request) {
}
