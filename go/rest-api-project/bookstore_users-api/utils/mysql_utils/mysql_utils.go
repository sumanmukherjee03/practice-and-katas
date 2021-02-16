package mysql_utils

import (
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/errors"
)

const (
	SqlErrorNoRows = "no rows in result set"
)

func ParseError(err error) *errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), SqlErrorNoRows) {
			return errors.NewNotFoundError(fmt.Errorf("Could not find record with the given id - %v", err))
		}
		return errors.NewInternalServerError(fmt.Errorf("Encountered an internal server error when parsing database error - %v", err))
	}
	switch sqlErr.Number {
	case 1062: // Represents a Duplicate Key error when inserting in a database
		return errors.NewBadRequestError(fmt.Errorf("Key for record in insert already exists - %v", err))
	}
	return errors.NewInternalServerError(fmt.Errorf("Encountered an internal server error when persisting to the database - %v", err))
}
