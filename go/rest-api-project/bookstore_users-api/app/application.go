package app

import "github.com/gin-gonic/gin"

var (
	// Remember that each request is handled by a different goroutine
	// So, keep the handlers free of any shared state
	router = gin.Default()
)

func StartApplication() {
	mapUrls()
	if err := router.Run(":8080"); err != nil {
		panic("Server failed to start")
	}
}
