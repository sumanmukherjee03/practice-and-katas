package main

import (
	"net/http"
)

func (app *application) getOneGenre(w http.ResponseWriter, r *http.Request) {
	id, err := app.getIdFromUrlParams(w, r)
	if err != nil {
		app.badRequestErrorJSON(w, err)
		return
	}
	genre, err := app.models.DB.GetGenreByID(id)
	if err != nil {
		app.entityNotFoundErrorJSON(w, err)
		return
	}
	if err = app.writeJSON(w, http.StatusOK, genre, "genre"); err != nil {
		app.serverErrorJSON(w, err)
		return
	}
}

func (app *application) getAllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := app.models.DB.GetAllGenres()
	if err != nil {
		app.serverErrorJSON(w, err)
		return
	}
	if err = app.writeJSON(w, http.StatusOK, genres, "genres"); err != nil {
		app.serverErrorJSON(w, err)
		return
	}
}

func (app *application) getAllMoviesByGenre(w http.ResponseWriter, r *http.Request) {
	genreID, err := app.getIdFromUrlParams(w, r)
	if err != nil {
		app.badRequestErrorJSON(w, err)
		return
	}
	movies, err := app.models.DB.GetAllMovies(genreID)
	if err != nil {
		app.serverErrorJSON(w, err)
		return
	}
	if err = app.writeJSON(w, http.StatusOK, movies, "movies"); err != nil {
		app.serverErrorJSON(w, err)
		return
	}
}
