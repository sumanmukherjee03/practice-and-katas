package rest_errors

import (
	"errors"
	"fmt"
	"net/http"
)

type RestErr struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func NewError(msg string) error {
	return errors.New(msg)
}

func NewBadRequestError(err error) *RestErr {
	return &RestErr{
		Message: fmt.Sprintf("Encountered a bad request - %v", err.Error()),
		Status:  http.StatusBadRequest,
		Error:   "bad_request",
	}
}

func NewUnauthorizedError(err error) *RestErr {
	return &RestErr{
		Message: fmt.Sprintf("Invalid access token - %v", err.Error()),
		Status:  http.StatusUnauthorized,
		Error:   "unauthorized",
	}
}

func NewNotFoundError(err error) *RestErr {
	return &RestErr{
		Message: fmt.Sprintf("Entity could not be found - %v", err.Error()),
		Status:  http.StatusNotFound,
		Error:   "not_found",
	}
}

func NewInternalServerError(err error) *RestErr {
	return &RestErr{
		Message: fmt.Sprintf("Encountered an internal server error - %v", err.Error()),
		Status:  http.StatusInternalServerError,
		Error:   "internal_server_error",
	}
}
