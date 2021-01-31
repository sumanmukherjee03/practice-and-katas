package services

import (
	"fmt"
	"net/http"

	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/domain/users"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/errors"
)

// By convention, always return error at the end
func CreateUser(user users.User) (*users.User, *errors.RestErr) {
	var err error
	if err != nil {
		restErr := errors.RestErr{
			Message: fmt.Sprintf("Encountered an internal server error when processing request - %v", err.Error()),
			Status:  http.StatusInternalServerError,
			Error:   "internal_server_error",
		}
		return nil, &restErr
	}
	return &user, nil
}
