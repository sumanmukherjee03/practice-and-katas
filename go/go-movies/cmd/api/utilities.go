package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, wrap string) error {
	wrapper := make(map[string]interface{})
	wrapper[wrap] = data
	jsonData, err := json.Marshal(wrapper)
	if err != nil {
		return fmt.Errorf("ERROR : Could not marshal json responding to http request - %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status) // This sends a http response header with the proper status code
	w.Write(jsonData)
	return nil
}
