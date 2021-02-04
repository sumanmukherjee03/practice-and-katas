package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/domain/users"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/services"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/errors"
)

//////////////////////// HELPER FUNCS ////////////////////
func getUserId(userIdStr string) (int64, *errors.RestErr) {
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		restErr := errors.NewBadRequestError(err)
		return 0, restErr
	}
	return userId, nil
}

//////////////////////// PUBLIC CONTROLLER FUNCS ////////////////////

func Create(ctx *gin.Context) {
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

func Get(ctx *gin.Context) {
	userId, userIdRestErr := getUserId(ctx.Param("user_id"))
	if userIdRestErr != nil {
		ctx.JSON(userIdRestErr.Status, userIdRestErr)
		return
	}

	user, notFoundErr := services.GetUser(userId)
	if notFoundErr != nil {
		ctx.JSON(notFoundErr.Status, notFoundErr)
		return
	}

	ctx.JSON(http.StatusFound, user)
}

func Update(ctx *gin.Context) {
	userId, userIdRestErr := getUserId(ctx.Param("user_id"))
	if userIdRestErr != nil {
		ctx.JSON(userIdRestErr.Status, userIdRestErr)
		return
	}

	var user users.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError(err)
		ctx.JSON(restErr.Status, restErr)
		return
	}

	user.Id = userId

	isPartialUpdate := ctx.Request.Method == http.MethodPatch

	res, serverErr := services.UpdateUser(isPartialUpdate, user)
	if serverErr != nil {
		ctx.JSON(serverErr.Status, serverErr)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func Delete(ctx *gin.Context) {
	userId, userIdRestErr := getUserId(ctx.Param("user_id"))
	if userIdRestErr != nil {
		ctx.JSON(userIdRestErr.Status, userIdRestErr)
		return
	}

	if serverErr := services.DeleteUser(userId); serverErr != nil {
		ctx.JSON(serverErr.Status, serverErr)
		return
	}
	ctx.JSON(http.StatusOK, userId)
}
