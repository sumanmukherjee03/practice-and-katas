package users

import (
	"fmt"
	"strings"

	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/datasources/mysql/usersdb"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/mysql_utils"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_utils-go/logger"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_utils-go/rest_errors"
)

var (
	usersDB = make(map[int64]*User)
	log     = logger.GetLogger()
)

const (
	queryInsertUser                 = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser                    = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id=?;"
	queryUpdateUser                 = "UPDATE users SET first_name=?, last_name=?, email=?, status=?, password=? WHERE id=?;"
	queryDeleteUser                 = "DELETE FROM users WHERE id = ?;"
	queryFindUsersByStatus          = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
	queryFindUserByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE email=? AND password=? AND status=?;"
)

// A code chunk if you simply want to ping the database to test a connection
// if err := usersdb.Client.Ping(); err != nil {
// panic(err)
// }

func (u *User) Get() *rest_errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryGetUser) // Prepare a DB statement first. Prepared DB statements are also more performant.
	if err != nil {
		log.Error("Error in preparing statement for fetching user from database", err)
		return rest_errors.NewInternalServerError(rest_errors.NewError("database error"))
	}
	defer stmt.Close() // Make sure you defer close the statement to not have idle connections lingering around
	// stmt.QueryRow returns a single row and the connection closes automatically on return
	// However, if we used stmt.Query, it would have returned *Rows in which case, we would have
	// had a need to defer rows.Close() so that we dont end up with idle connections to the DB.
	if getErr := stmt.QueryRow(u.Id).Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.DateCreated, &u.Status); getErr != nil {
		log.Error("Error in trying to get user by id from database", getErr)
		return rest_errors.NewInternalServerError(rest_errors.NewError("database error"))
		// return mysql_utils.ParseError(getErr)
	}
	return nil
}

func (u *User) Save() *rest_errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryInsertUser) // Prepare a DB statement first. Prepared DB statements are also more performant.
	if err != nil {
		log.Error("Error in preparing statement for saving user into database", err)
		return rest_errors.NewInternalServerError(rest_errors.NewError("database error"))
	}
	defer stmt.Close() // Make sure you defer close the statement to not have idle connections lingering around
	insertRes, insertErr := stmt.Exec(u.FirstName, u.LastName, u.Email, u.DateCreated, u.Status, u.Password)
	if insertErr != nil {
		log.Error("Error in inserting user into database", insertErr)
		return rest_errors.NewInternalServerError(rest_errors.NewError("database error"))
		// return mysql_utils.ParseError(insertErr)
	}
	userId, lastInsertIdErr := insertRes.LastInsertId() // Get the id of the row just inserted
	if lastInsertIdErr != nil {
		log.Error("Error in getting id of last inserted user into database", lastInsertIdErr)
		return rest_errors.NewInternalServerError(rest_errors.NewError("database error"))
		// return rest_errors.NewInternalServerError(lastInsertIdErr)
	}
	u.Id = userId
	return nil
}

func (u *User) Update() *rest_errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryUpdateUser) // Prepare a DB statement first. Prepared DB statements are also more performant.
	if err != nil {
		log.Error("Error in preparing statement for updating user in database", err)
		return rest_errors.NewInternalServerError(rest_errors.NewError("database error"))
	}
	defer stmt.Close() // Make sure you defer close the statement to not have idle connections lingering around
	_, updateErr := stmt.Exec(u.FirstName, u.LastName, u.Email, u.Status, u.Password, u.Id)
	if updateErr != nil {
		log.Error("Error in updating user in database", updateErr)
		return rest_errors.NewInternalServerError(rest_errors.NewError("database error"))
		// return mysql_utils.ParseError(updateErr)
	}
	return nil
}

func (u *User) Delete() *rest_errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryDeleteUser) // Prepare a DB statement first. Prepared DB statements are also more performant.
	if err != nil {
		log.Error("Error in preparing statement for updating user in database", err)
		return rest_errors.NewInternalServerError(rest_errors.NewError("database error"))
	}
	defer stmt.Close() // Make sure you defer close the statement to not have idle connections lingering around
	_, deleteErr := stmt.Exec(u.Id)
	if deleteErr != nil {
		log.Error("Error in deleting user from database", deleteErr)
		return rest_errors.NewInternalServerError(rest_errors.NewError("database error"))
		// return mysql_utils.ParseError(deleteErr)
	}
	return nil
}

func (u *User) FindByEmailAndPassword() *rest_errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryFindUserByEmailAndPassword)
	if err != nil {
		log.Error("Error in preparing statement for finding user by email and password in database", err)
		return rest_errors.NewInternalServerError(rest_errors.NewError("database error"))
	}
	defer stmt.Close() // Make sure you defer close the statement to not have idle connections lingering around

	if getErr := stmt.QueryRow(u.Email, u.Password, StatusActive).Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.DateCreated, &u.Status); getErr != nil {
		log.Error("Error in trying to get user by email and password from database", getErr)
		if strings.Contains(getErr.Error(), mysql_utils.SqlErrorNoRows) {
			return rest_errors.NewNotFoundError(rest_errors.NewError("invalid user credentials"))
		}
		return rest_errors.NewInternalServerError(rest_errors.NewError("database error"))
		// return mysql_utils.ParseError(getErr)
	}

	return nil
}

func FindByStatus(status string) ([]User, *rest_errors.RestErr) {
	stmt, err := usersdb.Client.Prepare(queryFindUsersByStatus)
	if err != nil {
		log.Error("Error in preparing statement for finding user by status in database", err)
		return nil, rest_errors.NewInternalServerError(rest_errors.NewError("database error"))
	}
	defer stmt.Close() // Make sure you defer close the statement to not have idle connections lingering around

	rows, err := stmt.Query(status) // Use query instead of exec to get back rows of results
	if err != nil {
		log.Error("Error in running query to find users by status from database", err)
		return nil, rest_errors.NewInternalServerError(rest_errors.NewError("database error"))
		// return nil, rest_errors.NewInternalServerError(err)
	}
	defer rows.Close() // Make sure you defer close the rows to not have idle connections lingering around

	res := make([]User, 0)
	for rows.Next() { // keep iterating as long as there is a next row
		var user User
		if getErr := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
			log.Error("Error in getting single row from database when finding users by status", getErr)
			return nil, rest_errors.NewInternalServerError(rest_errors.NewError("database error"))
			// return nil, mysql_utils.ParseError(getErr)
		}
		res = append(res, user)
	}

	if len(res) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Errorf("No user matching status %s found", status))
	}

	return res, nil
}
