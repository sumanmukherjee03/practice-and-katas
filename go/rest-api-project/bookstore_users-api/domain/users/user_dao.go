package users

import (
	"fmt"

	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/datasources/mysql/usersdb"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/date_utils"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/errors"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/mysql_utils"
)

var (
	usersDB = make(map[int64]*User)
)

const (
	queryInsertUser        = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?, ?, ?, ?);"
	queryGetUser           = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id=?;"
	queryUpdateUser        = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser        = "DELETE FROM users WHERE id = ?;"
	queryFindUsersByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
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
		return mysql_utils.ParseError(getErr)
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
		return mysql_utils.ParseError(insertErr)
	}
	userId, lastInsertIdErr := insertRes.LastInsertId() // Get the id of the row just inserted
	if lastInsertIdErr != nil {
		return errors.NewInternalServerError(lastInsertIdErr)
	}
	u.Id = userId
	return nil
}

func (u *User) Update() *errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryUpdateUser) // Prepare a DB statement first. Prepared DB statements are also more performant.
	if err != nil {
		return errors.NewInternalServerError(err)
	}
	defer stmt.Close() // Make sure you defer close the statement to not have idle connections lingering around
	_, updateErr := stmt.Exec(u.FirstName, u.LastName, u.Email, u.Id)
	if updateErr != nil {
		return mysql_utils.ParseError(updateErr)
	}
	return nil
}

func (u *User) Delete() *errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryDeleteUser) // Prepare a DB statement first. Prepared DB statements are also more performant.
	if err != nil {
		return errors.NewInternalServerError(err)
	}
	defer stmt.Close() // Make sure you defer close the statement to not have idle connections lingering around
	_, deleteErr := stmt.Exec(u.Id)
	if deleteErr != nil {
		return mysql_utils.ParseError(deleteErr)
	}
	return nil
}

func FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := usersdb.Client.Prepare(queryFindUsersByStatus)
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	defer stmt.Close() // Make sure you defer close the statement to not have idle connections lingering around
	rows, err := stmt.Query(status)
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	defer rows.Close() // Make sure you defer close the rows to not have idle connections lingering around

	res := make([]User, 0)
	for rows.Next() {
		var user User
		if getErr := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
			return nil, mysql_utils.ParseError(getErr)
		}
		res = append(res, user)
	}

	if len(res) == 0 {
		return nil, errors.NewNotFoundError(fmt.Errorf("No user matching status %s found", status))
	}

	return res, nil
}
