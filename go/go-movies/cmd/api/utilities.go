package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type jsonError struct {
	Message string `json:"message"`
}

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

func (app *application) clientErrorJSON(w http.ResponseWriter, err error) {
	e := jsonError{
		Message: fmt.Sprintf("CLIENT ERROR : %v", err),
	}
	app.logger.Print(e)
	app.writeJSON(w, http.StatusBadRequest, e, "error")
}

func (app *application) serverErrorJSON(w http.ResponseWriter, err error) {
	e := jsonError{
		Message: fmt.Sprintf("SERVER ERROR : %v", err),
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
