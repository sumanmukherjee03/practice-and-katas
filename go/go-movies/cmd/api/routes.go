package main

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// This is a a wrapper function for http handlers which we will be using to wrap
// our middlewares so thatrequest params are available upstream in application http handlers
func (app *application) wrapMiddleware(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		// Add query parameters to the request context so that they are available upstream in the application http handlers
		ctx := context.WithValue(r.Context(), httprouter.ParamsKey, params)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// Remember that httprouter.Router satisfies the interface to be a server mux
// The signature of this function used to be - func (app *application) routes() *httprouter.Router
// But due to the addition of the CORS middleware we had to change the signature to this new version.
// The return type has changed to http.Handler. http.Handler is an interface on which you can call the ServeHTTP function.
// The router satisfies this interface as well. So, no need of change in the object that we are returning.
// But to enable CORS in every single route we return app.enableCORS(router).
// And since router satisfies the http.Handler interface this is possible.
func (app *application) routes() http.Handler {
	router := httprouter.New()
	secure := alice.New(app.authenticate)

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc(http.MethodPost, "/v1/graphql/movies", app.moviesGraphQL)

	router.HandlerFunc(http.MethodPost, "/v1/signin", app.signin)

	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getOneMovie)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getAllMovies)
	router.HandlerFunc(http.MethodGet, "/v1/genre/:id", app.getOneGenre)
	router.HandlerFunc(http.MethodGet, "/v1/genre/:id/movies", app.getAllMoviesByGenre)
	router.HandlerFunc(http.MethodGet, "/v1/genres", app.getAllGenres)

	// Previously the routes used to look like this
	//     router.HandlerFunc(http.MethodPost, "/v1/admin/movie/edit", app.editMovie)
	//     router.HandlerFunc(http.MethodDelete, "/v1/admin/movie/:id/delete", app.deleteMovie)
	// But we are changing the declaration of that same route slightly so that we can use the middleware chaining
	router.POST("/v1/admin/movie/edit", app.wrapMiddleware(secure.ThenFunc(app.editMovie)))
	router.DELETE("/v1/admin/movie/:id/delete", app.wrapMiddleware(secure.ThenFunc(app.deleteMovie)))

	return app.enableCORS(router)
}
