package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	version = "1.0.0"
)

type config struct {
	port int
	env  string
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "application development environment")
	flag.Parse()
	fmt.Println("running")

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		currentStatus := AppStatus{
			Status:      "available",
			Environment: cfg.env,
			Version:     version,
		}
		data, err := json.MarshalIndent(currentStatus, "", "  ")
		if err != nil {
			log.Error(fmt.Errorf("ERROR : Could not marshal json for reporting status - %v", err))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.port), nil)
	if err != nil {
		log.Error(fmt.Errorf("ERROR : Could not start server - %v", err))
	}
}
