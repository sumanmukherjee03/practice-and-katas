package app

import (
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/controllers/ping"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/controllers/users"
)

func mapUrls() {
	router.GET("/ping", ping.Ping)

	router.GET("/users/:user_id", users.GetUser)
	router.POST("/users", users.CreateUser)
	router.DELETE("/users/:user_id", users.DeleteUser)
}
