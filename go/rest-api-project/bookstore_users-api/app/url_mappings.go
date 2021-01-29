package app

import "github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/controllers"

func mapUrls() {
	router.GET("/ping", controllers.Ping)
}
