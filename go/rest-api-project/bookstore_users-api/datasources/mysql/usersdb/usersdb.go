package usersdb

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql" // This implements the sql interface in the database/sql package
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/env_utils"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_utils-go/logger"
)

const (
	mysql_username_env_var = "MYSQL_USERNAME"
	mysql_password_env_var = "MYSQL_PASSWORD"
	mysql_host_env_var     = "MYSQL_HOST"
	mysql_port_env_var     = "MYSQL_PORT"
	mysql_schema_env_var   = "MYSQL_SCHEMA"
)

var (
	Client *sql.DB
	log    = logger.GetLogger()
)

// init functions are initialization functions in go and are called only once when a package is imported
// irrespective of how many times it is being imported. It is a good place to initialize things like
// connection pools, DB connections etc.
// For more details : https://tutorialedge.net/golang/the-go-init-function/
func init() {
	var err error
	dataSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8",
		env_utils.GetEnvVarWithDefaults(mysql_username_env_var, "root"),
		os.Getenv(mysql_password_env_var),
		env_utils.GetEnvVarWithDefaults(mysql_host_env_var, "localhost"),
		env_utils.GetEnvVarWithDefaults(mysql_port_env_var, "3306"),
		env_utils.GetEnvVarWithDefaults(mysql_schema_env_var, "bookstore_users"),
	)
	Client, err = sql.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}
	if err = Client.Ping(); err != nil {
		panic(err)
	}
	mysql.SetLogger(log)
	log.Info("Database successfully connected")
}
