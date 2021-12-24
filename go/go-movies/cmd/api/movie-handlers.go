package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go-movies/models"

	"github.com/julienschmidt/httprouter"
)

func (app *application) getOneMovie(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id <= 0 {
		app.logger.Print(fmt.Errorf("ERROR : Could not find a valid id in the url params - %v", err))
	}
	app.logger.Println("Id is", id)
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
	err = app.writeJSON(w, http.StatusOK, movie, "movie")
	if err != nil {
		app.logger.Print(err)
		return
	}
}

func (app *application) getAllMovies(w http.ResponseWriter, r *http.Request) {
}
