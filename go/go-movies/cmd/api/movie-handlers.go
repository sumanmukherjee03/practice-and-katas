package main

import (
	"net/http"
)

func (app *application) getOneMovie(w http.ResponseWriter, r *http.Request) {
	id, err := app.getIdFromUrlParams(w, r)
	if err != nil {
		app.badRequestErrorJSON(w, err)
		return
	}
	movie, err := app.models.DB.Get(id)
	if err != nil {
		app.entityNotFoundErrorJSON(w, err)
		return
	}
	if err = app.writeJSON(w, http.StatusOK, movie, "movie"); err != nil {
		app.serverErrorJSON(w, err)
		return
	}
}

func (app *application) getAllMovies(w http.ResponseWriter, r *http.Request) {
}
