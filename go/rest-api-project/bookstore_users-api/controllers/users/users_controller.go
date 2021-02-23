package users

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-go/oauth"
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

func isReqPublic(ctx *gin.Context) bool {
	return strings.Compare(ctx.GetHeader("X-Public"), "true") == 0
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

	createdUser, serverErr := services.UsersService.CreateUser(user)
	if serverErr != nil {
		ctx.JSON(serverErr.Status, serverErr)
		return
	}

	ctx.JSON(http.StatusCreated, createdUser.Marshal(isReqPublic(ctx)))
}

func Get(ctx *gin.Context) {
	if err := oauth.Authenticate(ctx.Request); err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	userId, userIdRestErr := getUserId(ctx.Param("user_id"))
	if userIdRestErr != nil {
		ctx.JSON(userIdRestErr.Status, userIdRestErr)
		return
	}

	user, notFoundErr := services.UsersService.GetUser(userId)
	if notFoundErr != nil {
		ctx.JSON(notFoundErr.Status, notFoundErr)
		return
	}

	// Check if user is asking for their own information.
	// If yes, then make a private request with all the details of the user.
	if oauth.GetCallerId(ctx.Request) == user.Id {
		ctx.JSON(http.StatusOK, user.Marshal(false))
		return
	}

	ctx.JSON(http.StatusOK, user.Marshal(oauth.IsPublic(ctx.Request)))
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

	updatedUser, serverErr := services.UsersService.UpdateUser(isPartialUpdate, user)
	if serverErr != nil {
		ctx.JSON(serverErr.Status, serverErr)
		return
	}

	ctx.JSON(http.StatusOK, updatedUser.Marshal(isReqPublic(ctx)))
}

func Delete(ctx *gin.Context) {
	userId, userIdRestErr := getUserId(ctx.Param("user_id"))
	if userIdRestErr != nil {
		ctx.JSON(userIdRestErr.Status, userIdRestErr)
		return
	}

	if serverErr := services.UsersService.DeleteUser(userId); serverErr != nil {
		ctx.JSON(serverErr.Status, serverErr)
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{"status": fmt.Sprintf("Deleted user %d", userId)})
}

func Login(ctx *gin.Context) {
	var loginReq users.LoginRequest
	if err := ctx.ShouldBindJSON(&loginReq); err != nil {
		restErr := errors.NewBadRequestError(err)
		ctx.JSON(restErr.Status, restErr)
		return
	}
	loginReq.PrepBeforeSubmit()

	user, notFoundErr := services.UsersService.LoginUser(loginReq)
	if notFoundErr != nil {
		ctx.JSON(notFoundErr.Status, notFoundErr)
		return
	}

	ctx.JSON(http.StatusOK, user.Marshal(isReqPublic(ctx)))
}

func Search(ctx *gin.Context) {
	status := ctx.Query("status") // Since status is coming as a query parameter and not as a paramter in the url

	searchedUsers, serverErr := services.UsersService.SearchUser(status)
	if serverErr != nil {
		ctx.JSON(serverErr.Status, serverErr)
		return
	}

	ctx.JSON(http.StatusOK, searchedUsers.Marshal(isReqPublic(ctx)))
}
