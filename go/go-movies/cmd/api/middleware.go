package main

import "net/http"

// All http middlewares in go commonly have this signature where it takes
// a next http.Handler and returns a http.Handler
func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins to connect to the backend
		next.ServeHTTP(w, r)
	})
}
