package main

import (
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
	if err := app.writeJSON(w, http.StatusOK, currentStatus, "health"); err != nil {
		app.serverErrorJSON(w, err)
		return
	}
}
