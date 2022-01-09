package main

import (
	"fmt"
	"go-movies/models"
	"io"
	"net/http"
	"strings"

	"github.com/graphql-go/graphql"
)

var (
	movies []*models.Movie

	// Graphql schema definition
	fields = graphql.Fields{
		"movie": &graphql.Field{
			Type:        movieType,
			Description: "Get movie by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					for _, movie := range movies {
						if movie.ID == id {
							return movie, nil
						}
					}
				}
				return nil, nil
			},
		},
		"list": &graphql.Field{
			Type:        graphql.NewList(movieType),
			Description: "Get all movies",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return movies, nil
			},
		},
		"search": &graphql.Field{
			Type:        graphql.NewList(movieType),
			Description: "Search movies by title",
			Args: graphql.FieldConfigArgument{
				"titleContains": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var res []*models.Movie
				search, ok := p.Args["titleContains"].(string)
				if !ok {
					return res, fmt.Errorf("Could not get the param `titleContains` from the graphql call")
				}
				for _, currentMovie := range movies {
					if strings.Contains(currentMovie.Title, search) {
						res = append(res, currentMovie)
					}
				}
				return res, nil
			},
		},
	}

	// This represents the kind of information we would expose to the client from the database
	movieType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Movie",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"title": &graphql.Field{
					Type: graphql.String,
				},
				"description": &graphql.Field{
					Type: graphql.String,
				},
				"year": &graphql.Field{
					Type: graphql.Int,
				},
				"release_date": &graphql.Field{
					Type: graphql.DateTime,
				},
				"runtime": &graphql.Field{
					Type: graphql.Int,
				},
				"rating": &graphql.Field{
					Type: graphql.Int,
				},
				"mpaa_rating": &graphql.Field{
					Type: graphql.Int,
				},
				"created_at": &graphql.Field{
					Type: graphql.DateTime,
				},
				"updated_at": &graphql.Field{
					Type: graphql.DateTime,
				},
			},
		},
	)
)

func (app *application) moviesGraphQL(w http.ResponseWriter, r *http.Request) {
	var err error

	// Populate our package level variable called movies here. See the reference to it in the Resolve functions of graphql.Field
	movies, err = app.models.DB.GetAllMovies()
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

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields} // fields = the variable declared at top representing the graphql schema
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		app.serverErrorJSON(w, fmt.Errorf("Couldnt get a graphql schema from the schema config - %v", err))
		return
	}

	params := graphql.Params{Schema: schema, RequestString: query}
	resp := graphql.Do(params) // This is when the Resolve defined in the graphql.Field callback above kicks in
	if len(resp.Errors) > 0 {
		app.serverErrorJSON(w, fmt.Errorf("Couldnt get a graphql response from the backend - %v", resp.Errors))
		return
	}

	if err = app.writeJSON(w, http.StatusOK, resp, ""); err != nil {
		app.serverErrorJSON(w, err)
		return
	}
}
