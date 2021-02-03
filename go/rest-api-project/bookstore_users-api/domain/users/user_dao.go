package users

import (
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/datasources/mysql/usersdb"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/date_utils"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/errors"
)

var (
	usersDB = make(map[int64]*User)
)

const (
	sqlErrorNoRows  = "no rows in result set"
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?, ?, ?, ?);"
	queryGetUser    = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id=?;"
)

// A code chunk if you simply want to ping the database to test a connection
// if err := usersdb.Client.Ping(); err != nil {
// panic(err)
// }

func (u *User) Get() *errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryGetUser) // Prepare a DB statement first. Prepared DB statements are also more performant.
	if err != nil {
		return errors.NewInternalServerError(err)
	}
	defer stmt.Close() // Make sure you defer close the statement to not have idle connections lingering around
	// stmt.QueryRow returns a single row and the connection closes automatically on return
	// However, if we used stmt.Query, it would have returned *Rows in which case, we would have
	// had a need to defer rows.Close() so that we dont end up with idle connections to the DB.
	if getErr := stmt.QueryRow(u.Id).Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.DateCreated); getErr != nil {
		if strings.Contains(getErr.Error(), sqlErrorNoRows) {
			return errors.NewNotFoundError(fmt.Errorf("Could not find an user with the given id : %d", u.Id))
		}
		return errors.NewInternalServerError(getErr)
	}
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
		sqlErr, ok := insertErr.(*mysql.MySQLError)
		if !ok {
			return errors.NewInternalServerError(fmt.Errorf("Could not convert error returned by DB into a mysql error : %v", insertErr))
		}
		switch sqlErr.Number {
		case 1062: // Represents a Duplicate Key error when inserting in a database
			return errors.NewBadRequestError(fmt.Errorf("Email for user already exists : %s", u.Email))
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
