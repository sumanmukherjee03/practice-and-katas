package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/domain/users"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/services"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/errors"
)

func CreateUser(ctx *gin.Context) {
	var user users.User

	// The lines below can be replaced by the ctx.ShouldBindJSON function call
	// bytes, err := ioutil.ReadAll(ctx.Request.Body)
	// if err != nil {
	// fmt.Println(err)
	// return
	// }
	// if err := json.Unmarshal(bytes, &user); err != nil {
	// fmt.Println(err)
	// return
	// }

	// ctx.ShouldBindJSON does the job of receiving bytes array from the request body in POST
	// unmarshall it and populate the user struct
	if err := ctx.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError(err)
		ctx.JSON(restErr.Status, restErr)
		return
	}

	res, serverErr := services.CreateUser(user)
	if serverErr != nil {
		ctx.JSON(serverErr.Status, serverErr)
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func GetUser(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "GetUser NOT IMPLEMENTED")
}

func DeleteUser(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "DeleteUser NOT IMPLEMENTED")
}
