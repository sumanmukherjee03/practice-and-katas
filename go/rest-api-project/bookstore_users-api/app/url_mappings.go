package app

import (
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/controllers/ping"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/controllers/users"
)

func mapUrls() {
	router.GET("/ping", ping.Ping)

	router.GET("/users/:user_id", users.Get)
	router.POST("/users", users.Create)
	router.PUT("/users/:user_id", users.Update)
	router.PATCH("/users/:user_id", users.Update)
	router.DELETE("/users/:user_id", users.Delete)
	router.GET("/internal/users/search", users.Search)
}
