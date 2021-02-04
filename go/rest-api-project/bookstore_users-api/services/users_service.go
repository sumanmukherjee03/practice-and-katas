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

func GetUser(userId int64) (*users.User, *errors.RestErr) {
	var u = &users.User{Id: userId}
	if err := u.Get(); err != nil {
		return nil, err
	}
	return u, nil
}

// By convention, always return error at the end
func UpdateUser(isPartial bool, u users.User) (*users.User, *errors.RestErr) {
	var err *errors.RestErr
	currentUser, err := GetUser(u.Id)
	if err != nil {
		return nil, err // Check if user even exists in DB and return an error if it doesnt
	}

	// Handle Patch and Put type methods
	if isPartial {
		if len(u.FirstName) > 0 {
			currentUser.FirstName = u.FirstName
		}
		if len(u.LastName) > 0 {
			currentUser.LastName = u.LastName
		}
		if len(u.Email) > 0 {
			currentUser.Email = u.Email
		}
	} else {
		currentUser.FirstName = u.FirstName
		currentUser.LastName = u.LastName
		currentUser.Email = u.Email
	}

	// Validate that the user being passed in is valid
	if err = currentUser.Validate(); err != nil {
		return nil, err
	}

	if err = currentUser.Update(); err != nil {
		return nil, err
	}

	return currentUser, nil
}

func DeleteUser(userId int64) *errors.RestErr {
	var u = &users.User{Id: userId}
	if err := u.Delete(); err != nil {
		return err
	}
	return nil
}
