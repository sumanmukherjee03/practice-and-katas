package users

import (
	"fmt"

	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/errors"
)

var (
	usersDB = make(map[int64]*User)
)

func (u *User) Get() *errors.RestErr {
	res, found := usersDB[u.Id]
	if !found {
		return errors.NewNotFoundError(fmt.Errorf("User with id %d not found", u.Id))
	}
	u.Id = res.Id
	u.FirstName = res.FirstName
	u.LastName = res.LastName
	u.Email = res.Email
	u.DateCreated = res.DateCreated
	return nil
}

func (u *User) Save() *errors.RestErr {
	if res, found := usersDB[u.Id]; found {
		if res.Email == u.Email {
			return errors.NewBadRequestError(fmt.Errorf("User with email %s already exists", u.Email))
		}
		return errors.NewBadRequestError(fmt.Errorf("User with id %d already exists", u.Id))
	}
	usersDB[u.Id] = u
	return nil
}