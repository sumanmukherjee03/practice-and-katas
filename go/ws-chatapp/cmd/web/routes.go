package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/sumanmukherjee03/practice-and-katas/go/ws-chatapp/internal/handlers"
)

func routes() http.Handler {
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(handlers.Home))
	mux.Get("/ws", http.HandlerFunc(handlers.WsEndpoint))

	// Make sure to add the trailing slash at the end of the directory name
	fileServer := http.FileServer(http.Dir("./static/"))
	// Strip out the /static from the path and pass it to the fileserver
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
