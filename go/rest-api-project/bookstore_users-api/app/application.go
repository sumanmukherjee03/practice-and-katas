package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/logger"
)

var (
	// Remember that each request is handled by a different goroutine
	// So, keep the handlers free of any shared state
	router = gin.Default()
	log    = logger.GetLogger()
)

func StartApplication() {
	mapUrls()
	log.Info("Starting bookstore users application")
	if err := router.Run(":8080"); err != nil {
		panic("Server failed to start")
	}
}
