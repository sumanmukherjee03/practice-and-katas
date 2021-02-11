package services

import (
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/domain/users"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/errors"
)

var (
	UsersService usersServiceInterface = &usersService{} // Use the UsersService object as a singleton
)

type usersService struct{}

// We create this interface so that it is easy to mock in tests
type usersServiceInterface interface {
	CreateUser(users.User) (*users.User, *errors.RestErr)
	GetUser(int64) (*users.User, *errors.RestErr)
	UpdateUser(bool, users.User) (*users.User, *errors.RestErr)
	DeleteUser(int64) *errors.RestErr
	SearchUser(string) (users.Users, *errors.RestErr)
}

// By convention, always return error at the end
func (s *usersService) CreateUser(u users.User) (*users.User, *errors.RestErr) {
	var err *errors.RestErr
	u.PrepBeforeSave()
	if err = u.Validate(); err != nil {
		return nil, err
	}
	if err = u.Save(); err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *usersService) GetUser(userId int64) (*users.User, *errors.RestErr) {
	var u = &users.User{Id: userId}
	if err := u.Get(); err != nil {
		return nil, err
	}
	return u, nil
}

// By convention, always return error at the end
func (s *usersService) UpdateUser(isPartial bool, u users.User) (*users.User, *errors.RestErr) {
	var err *errors.RestErr
	currentUser, err := s.GetUser(u.Id)
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
		if len(u.Status) > 0 {
			currentUser.Status = u.Status
		}
		if len(u.Password) > 0 {
			currentUser.Password = u.Password
		}
	} else {
		currentUser.FirstName = u.FirstName
		currentUser.LastName = u.LastName
		currentUser.Email = u.Email
		currentUser.Status = u.Status
		currentUser.Password = u.Password
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

func (s *usersService) DeleteUser(userId int64) *errors.RestErr {
	var u = &users.User{Id: userId}
	return u.Delete()
}

func (s *usersService) SearchUser(status string) (users.Users, *errors.RestErr) {
	return users.FindByStatus(status) // Although users.FindByStatus returns a []User, it is treated as type Users
}
