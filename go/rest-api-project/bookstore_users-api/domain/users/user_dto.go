package users

import (
	"fmt"
	"strings"

	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/errors"
)

// User struct represents the Data Tranfer Object for the user entity
// Dont use a space after the json: for marshalling and unmarshalling to work properly
type User struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
}

// Validate method on user checks if an user is valid or not
func (u *User) Validate() *errors.RestErr {
	u.FirstName = strings.TrimSpace(strings.ToLower(u.FirstName))
	u.LastName = strings.TrimSpace(strings.ToLower(u.LastName))
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	if len(u.Email) == 0 {
		return errors.NewBadRequestError(fmt.Errorf("Invalid email address provided for user"))
	}
	return nil
}
