package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateUser(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "CreateUser NOT IMPLEMENTED")
}

func GetUser(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "GetUser NOT IMPLEMENTED")
}

func DeleteUser(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "DeleteUser NOT IMPLEMENTED")
}
