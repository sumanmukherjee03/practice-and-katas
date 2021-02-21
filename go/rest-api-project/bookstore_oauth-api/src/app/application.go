package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/clients/cassandra"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/http"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/repository/db"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/repository/rest"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/services/access_token"
)

var (
	router = gin.Default()
)

func StartApplication() {
	session := cassandra.GetSession()
	defer session.Close()
	atService := access_token.NewService(db.NewRepository(), rest.NewUsersRepository())
	atHandler := http.NewHandler(atService)
	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)
	if err := router.Run(":8080"); err != nil {
		panic("Server failed to start")
	}
}
