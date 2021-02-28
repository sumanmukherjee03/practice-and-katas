package http_utils

import (
	"encoding/json"
	"net/http"

	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_utils-go/rest_errors"
)

func RespondJson(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func RespondError(w http.ResponseWriter, err *rest_errors.RestErr) {
	RespondJson(w, err.Status, err)
}
