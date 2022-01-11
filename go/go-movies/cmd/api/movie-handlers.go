package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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

	// Only get the poster url if there isn't one
	if len(movie.Poster) == 0 {
		app.getPoster(&movie)
	}

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

////////////////////////////// HELPER FUNCS ////////////////////////////////////

func (app *application) getPoster(movie *models.Movie) {
	// type to receive the response of themoviedb api call
	// NOTE : You can get this if you copy the result of an API call from themoviedb.org into https://mholt.github.io/json-to-go/
	// And the api call can be done from the browser using : https://api.themoviedb.org/3/search/movie?api_key=<key>&query=The%20Shawshank%20Redemption
	type TheMovieDBResp struct {
		Page    int `json:"page"`
		Results []struct {
			Adult            bool    `json:"adult"`
			BackdropPath     string  `json:"backdrop_path"`
			GenreIds         []int   `json:"genre_ids"`
			ID               int     `json:"id"`
			OriginalLanguage string  `json:"original_language"`
			OriginalTitle    string  `json:"original_title"`
			Overview         string  `json:"overview"`
			Popularity       float64 `json:"popularity"`
			PosterPath       string  `json:"poster_path"`
			ReleaseDate      string  `json:"release_date"`
			Title            string  `json:"title"`
			Video            bool    `json:"video"`
			VoteAverage      float64 `json:"vote_average"`
			VoteCount        int     `json:"vote_count"`
		} `json:"results"`
		TotalPages   int `json:"total_pages"`
		TotalResults int `json:"total_results"`
	}

	client := &http.Client{}
	apiURL := "https://api.themoviedb.org/3/search/movie?api_key="
	movieTitle := url.QueryEscape(movie.Title)
	req, err := http.NewRequest("GET", apiURL+app.config.themoviedb.apikey+"&query="+movieTitle, nil)
	if err != nil {
		// If you encounter an error then dont fail - simply return and carry on
		log.Println(fmt.Sprintf("ERROR : Encountered an error making a new request to themoviedb.org to get movie poster - %v", err))
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		// If you receive an error response then dont fail - simply return and carry on
		log.Println(fmt.Sprintf("ERROR : Encountered an error response from themoviedb.org while getting a movie poster - %v", err))
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		// If you receive an error response then dont fail - simply return and carry on
		log.Println(fmt.Sprintf("ERROR : Encountered an error reading the body of response received from themoviedb.org while getting a movie poster - %v", err))
		return
	}
	var respObj TheMovieDBResp
	err = json.Unmarshal(bodyBytes, &respObj)
	if err != nil {
		// If you receive an error response then dont fail - simply return and carry on
		log.Println(fmt.Sprintf("ERROR : Encountered an error unmarshaling the body of response received from themoviedb.org while getting a movie poster - %v", err))
		return
	}

	if len(respObj.Results) > 0 {
		movie.Poster = respObj.Results[0].PosterPath
	}
}
