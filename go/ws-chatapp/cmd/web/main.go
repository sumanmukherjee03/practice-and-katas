package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	mux := routes()
	log.Info("Starting web server on port 8080")
	_ = http.ListenAndServe(":8080", mux)
}
