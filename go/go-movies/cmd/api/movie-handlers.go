package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go-movies/models"
)

// Dedidated type to handle marshaling and unmarshaling of movie payloads coming from the frontend as json
type MoviePayload struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Year        string `json:"year"`
	ReleaseDate string `json:"release_date"`
	Runtime     string `json:"runtime"`
	Rating      string `json:"rating"`
	MPAARating  string `json:"mpaa_rating"`
}

type DeleteMovieResponseJson struct {
	ID      int    `json:"id"`
	Status  string `json:"Status"`
	Message string `json:"message"`
}

func (app *application) getOneMovie(w http.ResponseWriter, r *http.Request) {
	id, err := app.getIdFromUrlParams(w, r)
	if err != nil {
		app.badRequestErrorJSON(w, err)
		return
	}
	movie, err := app.models.DB.GetMovieByID(id)
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
	movies, err := app.models.DB.GetAllMovies()
	if err != nil {
		app.serverErrorJSON(w, err)
		return
	}
	if err = app.writeJSON(w, http.StatusOK, movies, "movies"); err != nil {
		app.serverErrorJSON(w, err)
		return
	}
}

func (app *application) editMovie(w http.ResponseWriter, r *http.Request) {
	var payload MoviePayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		app.badRequestErrorJSON(w, err)
		return
	}

	var movie models.Movie

	if len(payload.ID) > 0 {
		movieID, err := strconv.Atoi(payload.ID)
		if err != nil {
			app.badRequestErrorJSON(w, err)
			return
		}
		if movieID > 0 {
			m, err := app.models.DB.GetMovieByID(movieID)
			if err != nil {
				app.entityNotFoundErrorJSON(w, err)
				return
			}
			movie = *m
		}
	}

	if len(payload.Runtime) > 0 {
		rt, err := strconv.Atoi(payload.Runtime)
		if err != nil {
			app.badRequestErrorJSON(w, err)
			return
		}
		movie.Runtime = rt
	}

	if len(payload.Rating) > 0 {
		r, err := strconv.Atoi(payload.Rating)
		if err != nil {
			app.badRequestErrorJSON(w, err)
			return
		}
		movie.Rating = r
	}

	if len(payload.ReleaseDate) > 0 {
		rd, err := time.Parse("2006-01-02", payload.ReleaseDate)
		if err != nil {
			app.badRequestErrorJSON(w, err)
			return
		}
		movie.ReleaseDate = rd
		movie.Year = movie.ReleaseDate.Year()
	}

	movie.Title = payload.Title
	movie.Description = payload.Description
	movie.MPAARating = payload.MPAARating

	if movie.ID == 0 {
		newID, err := app.models.DB.InsertMovie(movie)
		if err != nil {
			app.serverErrorJSON(w, err)
			return
		}
		movie.ID = newID
	} else {
		err := app.models.DB.UpdateMovie(movie)
		if err != nil {
			app.serverErrorJSON(w, err)
			return
		}
	}

	if err := app.writeJSON(w, http.StatusOK, movie, "movie"); err != nil {
		app.serverErrorJSON(w, err)
		return
	}
}

func (app *application) deleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := app.getIdFromUrlParams(w, r)
	if err != nil {
		app.badRequestErrorJSON(w, err)
		return
	}
	movie, err := app.models.DB.GetMovieByID(id)
	if err != nil {
		app.entityNotFoundErrorJSON(w, err)
		return
	}
	err = app.models.DB.DeleteMovie(movie.ID)
	if err != nil {
		app.serverErrorJSON(w, err)
		return
	}
	resp := DeleteMovieResponseJson{
		ID:      id,
		Status:  "OK",
		Message: "Successfully deleted the movie from the database",
	}
	if err = app.writeJSON(w, http.StatusOK, resp, "delete_response"); err != nil {
		app.serverErrorJSON(w, err)
		return
	}
}

func (app *application) searchMovies(w http.ResponseWriter, r *http.Request) {
}
