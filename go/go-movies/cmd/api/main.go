package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go-movies/models"

	_ "github.com/lib/pq"
)

const (
	version = "1.0.0"
)

// Type that stores application config
type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

// We are gonna use an instance of this type as a receiver in various other parts of our application
// This is an easy approach to share bit and pieces of functionality between various parts of our application
type application struct {
	config config
	logger *log.Logger
	models models.Models
}

func main() {
	var cfg config
	var dbUser, dbPassword, dbHost, dbName, dsn string

	flag.IntVar(&cfg.port, "port", 4000, "server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "application development environment")
	flag.StringVar(&dbUser, "db-user", "root", "application db user")
	flag.StringVar(&dbPassword, "db-password", "", "application db password")
	flag.StringVar(&dbHost, "db-host", "localhost", "application db host")
	flag.StringVar(&dbName, "db-name", "go_movies", "application db name")
	flag.Parse()

	dsn = fmt.Sprintf("postgres://%s@%s/%s?sslmode=disable", dbUser, dbHost, dbName)
	if len(dbPassword) > 0 {
		dsn = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbName)
	}
	cfg.db.dsn = dsn

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	app.logger.Println("Starting server on port", app.config.port)
	err = srv.ListenAndServe()
	if err != nil {
		app.logger.Println(fmt.Sprintf("ERROR : Could not start server - %v", err))
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
