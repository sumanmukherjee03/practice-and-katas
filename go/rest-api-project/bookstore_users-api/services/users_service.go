package services

import (
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/domain/users"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/errors"
)

// By convention, always return error at the end
func CreateUser(u users.User) (*users.User, *errors.RestErr) {
	var err *errors.RestErr
	if err = u.Validate(); err != nil {
		return nil, err
	}
	if err = u.Save(); err != nil {
		return nil, err
	}
	return &u, nil
}
