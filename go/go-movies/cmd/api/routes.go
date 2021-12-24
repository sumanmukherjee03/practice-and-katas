package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Remember that httprouter.Router satisfies the interface to be a server mux
func (app *application) routes() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
	return router
}
