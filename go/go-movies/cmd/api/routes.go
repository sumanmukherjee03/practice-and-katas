package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Remember that httprouter.Router satisfies the interface to be a server mux
// The signature of this function used to be - func (app *application) routes() *httprouter.Router
// But due to the addition of the CORS middleware we had to change the signature to this new version.
// The return type has changed to http.Handler. http.Handler is an interface on which you can call the ServeHTTP function.
// The router satisfies this interface as well. So, no need of change in the object that we are returning.
// But to enable CORS in every single route we return app.enableCORS(router).
// And since router satisfies the http.Handler interface this is possible.
func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getOneMovie)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getAllMovies)
	router.HandlerFunc(http.MethodGet, "/v1/genre/:id", app.getOneGenre)
	router.HandlerFunc(http.MethodGet, "/v1/genres", app.getAllGenres)
	return app.enableCORS(router)
}
