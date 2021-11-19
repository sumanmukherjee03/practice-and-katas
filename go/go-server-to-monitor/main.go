package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<h1>Success</h1>")
	})
	log.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", nil)
}
