package main

import (
	"fmt"
	"go-movies/models"
	"io"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
)

var (
	movies []*models.Movie
	fields = graphql.Fields{
		"movie": &graphql.Field{
			Type: movieType,
			Description: "Get movie by id",
			Args: ""
		}
	}
)

func (app *application) moviesGraphQL(w http.ResponseWriter, r *http.Request) {
	movies, err := app.models.DB.GetAllMovies()
	if err != nil {
		app.serverErrorJSON(w, fmt.Errorf("Couldnt fetch all movies from the DB - %v", err))
		return
	}

	q, err := io.ReadAll(r.Body)
	if err != nil {
		app.serverErrorJSON(w, fmt.Errorf("Couldnt read request body - %v", err))
		return
	}

	query := string(q)
	log.Println(query)
	log.Println(movies[0])
}
