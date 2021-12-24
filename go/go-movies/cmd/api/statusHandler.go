package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Output type for api call to fetch status
type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

func (app *application) statusHandler(w http.ResponseWriter, r *http.Request) {
	currentStatus := AppStatus{
		Status:      "available",
		Environment: app.config.env,
		Version:     version,
	}
	data, err := json.MarshalIndent(currentStatus, "", "  ")
	if err != nil {
		app.logger.Println(fmt.Sprintf("ERROR : Could not marshal json for reporting status - %v", err))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // This sends a http response header with the proper status code
	w.Write(data)
}
