package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/sumanmukherjee03/practice-and-katas/go/ws-chatapp/internal/handlers"
)

func routes() http.Handler {
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(handlers.Home))
	return mux
}
