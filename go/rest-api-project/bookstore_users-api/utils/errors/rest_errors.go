package errors

import (
	"fmt"
	"net/http"
)

type RestErr struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func NewBadRequestError(err error) *RestErr {
	return &RestErr{
		Message: fmt.Sprintf("Encountered a bad request - %v", err.Error()),
		Status:  http.StatusBadRequest,
		Error:   "bad_request",
	}
}

func NewNotFoundError(err error) *RestErr {
	return &RestErr{
		Message: fmt.Sprintf("Entity could not be found - %v", err.Error()),
		Status:  http.StatusNotFound,
		Error:   "not_found",
	}
}
