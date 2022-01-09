package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type jsonError struct {
	ErrorType string `json:"error_type"`
	Message   string `json:"message"`
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, wrap string) error {
	var jsonData []byte
	var err error
	if len(wrap) > 0 {
		wrapper := make(map[string]interface{})
		wrapper[wrap] = data
		jsonData, err = json.Marshal(wrapper)
	} else {
		jsonData, err = json.Marshal(data)
	}
	if err != nil {
		return fmt.Errorf("ERROR : Could not marshal json responding to http request - %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status) // This sends a http response header with the proper status code
	w.Write(jsonData)
	return nil
}

// errorJSON is a function to return a generic error based on a http status. It is a variadic function since it can take multiple status codes.
// This makes it an easy way to refactor code, since you can choose to pass 0 or more status codes when invoking this function
func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}
	e := jsonError{
		ErrorType: "ERROR",
		Message:   fmt.Sprintf("ERROR : %v", err),
	}
	app.logger.Print(e)
	app.writeJSON(w, statusCode, e, "error")
}

func (app *application) badRequestErrorJSON(w http.ResponseWriter, err error) {
	e := jsonError{
		ErrorType: "BAD_REQUEST_ERROR",
		Message:   fmt.Sprintf("ERROR : %v", err),
	}
	app.logger.Print(e)
	app.writeJSON(w, http.StatusBadRequest, e, "error")
}

func (app *application) authorizationErrorJSON(w http.ResponseWriter, err error) {
	e := jsonError{
		ErrorType: "AUTHORIZATION_ERROR",
		Message:   fmt.Sprintf("ERROR : %v", err),
	}
	app.logger.Print(e)
	app.writeJSON(w, http.StatusForbidden, e, "error")
}

func (app *application) entityNotFoundErrorJSON(w http.ResponseWriter, err error) {
	e := jsonError{
		ErrorType: "NOT_FOUND_ERROR",
		Message:   fmt.Sprintf("ERROR : %v", err),
	}
	app.logger.Print(e)
	app.writeJSON(w, http.StatusNotFound, e, "error")
}

func (app *application) serverErrorJSON(w http.ResponseWriter, err error) {
	e := jsonError{
		ErrorType: "SERVER_ERROR",
		Message:   fmt.Sprintf("ERROR : %v", err),
	}
	app.logger.Print(e)
	app.writeJSON(w, http.StatusInternalServerError, e, "error")
}

func (app *application) getIdFromUrlParams(w http.ResponseWriter, r *http.Request) (int, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id <= 0 {
		if err == nil {
			err = fmt.Errorf("Invalid id provided for lookup")
		}
		return 0, err
	}
	return id, nil
}
