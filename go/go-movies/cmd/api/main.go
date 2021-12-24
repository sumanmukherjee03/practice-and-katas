package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	version = "1.0.0"
)

// Type that stores application config
type config struct {
	port int
	env  string
}

// We are gonna use an instance of this type as a receiver in various other parts of our application
type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "application development environment")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	app.logger.Println("Starting server on port", app.config.port)
	err := srv.ListenAndServe()
	if err != nil {
		app.logger.Println(fmt.Sprintf("ERROR : Could not start server - %v", err))
	}
}
