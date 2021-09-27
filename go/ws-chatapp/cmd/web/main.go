package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/sumanmukherjee03/practice-and-katas/go/ws-chatapp/internal/handlers"
)

func main() {
	mux := routes()
	log.Info("Starting channel listener for websocket events")
	go handlers.ListenToWsChan()
	log.Info("Starting web server on port 8080")
	_ = http.ListenAndServe(":8080", mux)
}
