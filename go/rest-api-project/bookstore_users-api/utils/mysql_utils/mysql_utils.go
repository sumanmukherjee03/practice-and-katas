package mysql_utils

import (
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_utils-go/rest_errors"
)

const (
	SqlErrorNoRows = "no rows in result set"
)

func ParseError(err error) *rest_errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), SqlErrorNoRows) {
			return rest_errors.NewNotFoundError(fmt.Errorf("Could not find record with the given id - %v", err))
		}
		return rest_errors.NewInternalServerError(fmt.Errorf("Encountered an internal server error when parsing database error - %v", err))
	}
	switch sqlErr.Number {
	case 1062: // Represents a Duplicate Key error when inserting in a database
		return rest_errors.NewBadRequestError(fmt.Errorf("Key for record in insert already exists - %v", err))
	}
	return rest_errors.NewInternalServerError(fmt.Errorf("Encountered an internal server error when persisting to the database - %v", err))
}
