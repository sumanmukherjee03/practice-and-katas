package users

import (
	"fmt"
	"strings"

	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/crypto_utils"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/date_utils"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/errors"
)

const (
	StatusActive = "active"
)

// User struct represents the Data Tranfer Object for the user entity
// Dont use a space after the json: for marshalling and unmarshalling to work properly
type User struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	Password    string `json:password` // The value of json tag as -, states that ignore the password field in the struct when marshalling or unmarshalling
	// Password    string `json:-` // The value of json tag as -, states that ignore the password field in the struct when marshalling or unmarshalling
}

func (u *User) PrepBeforeSave() {
	u.FirstName = strings.TrimSpace(strings.ToLower(u.FirstName))
	u.LastName = strings.TrimSpace(strings.ToLower(u.LastName))
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	u.DateCreated = date_utils.GetNowDBFormat()
	u.Status = StatusActive
	u.Password = strings.TrimSpace(u.Password)
	u.Password = crypto_utils.GetMd5(u.Password)
}

// Validate method on user checks if an user is valid or not
func (u *User) Validate() *errors.RestErr {
	if len(u.Email) == 0 {
		return errors.NewBadRequestError(fmt.Errorf("Invalid email address provided for user"))
	}
	if len(u.Password) == 0 {
		return errors.NewBadRequestError(fmt.Errorf("Invalid password provided for user"))
	}
	if len(u.Status) == 0 {
		return errors.NewBadRequestError(fmt.Errorf("Invalid status provided for user"))
	}
	return nil
}
