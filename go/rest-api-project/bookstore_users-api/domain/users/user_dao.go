package users

import (
	"fmt"
	"strings"

	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/datasources/mysql/usersdb"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/date_utils"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/errors"
)

var (
	usersDB = make(map[int64]*User)
)

const (
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?, ?, ?, ?);"
)

func (u *User) Get() *errors.RestErr {
	if err := usersdb.Client.Ping(); err != nil {
		panic(err)
	}
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
	stmt, err := usersdb.Client.Prepare(queryInsertUser) // Prepare a DB statement first. Prepared DB statements are also more performant.
	if err != nil {
		return errors.NewInternalServerError(err)
	}
	defer stmt.Close() // Make sure you defer close the statement to not have idle connections lingering around
	u.DateCreated = date_utils.GetNowString()
	insertRes, insertErr := stmt.Exec(u.FirstName, u.LastName, u.Email, u.DateCreated)
	if insertErr != nil {
		if strings.Contains(insertErr.Error(), "email_unique") {
			return errors.NewBadRequestError(fmt.Errorf("Email already exists : %s", u.Email))
		}
		return errors.NewInternalServerError(insertErr)
	}
	userId, lastInsertIdErr := insertRes.LastInsertId() // Get the id of the row just inserted
	if lastInsertIdErr != nil {
		return errors.NewInternalServerError(lastInsertIdErr)
	}
	u.Id = userId
	return nil
}
