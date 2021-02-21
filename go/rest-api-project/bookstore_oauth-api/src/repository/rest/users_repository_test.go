package rest

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/domain/users"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/utils/errors"
)

// TestMain is the entrypoint of every test suite in a package is the TestMain function
// just like Main function is the entrypoint for a go package
func TestMain(m *testing.M) {
	setup()
	exitCode := m.Run()
	teardown()
	os.Exit(exitCode)
}

func setup() {
	httpmock.ActivateNonDefault(GetRestClient().GetClient())
}

func teardown() {
	httpmock.DeactivateAndReset()
}

func TestLoginUserSuccess(t *testing.T) {
	assert := assert.New(t)
	respUser := users.User{Id: 1, FirstName: "Foo", LastName: "Bar", Email: "foo.bar@baz.com"}
	responder := httpmock.NewJsonResponderOrPanic(http.StatusOK, respUser)
	httpmock.RegisterResponder("POST", "http://localhost:8081/users/login", responder)
	repo := usersRepository{}
	user, err := repo.LoginUser("foo.bar@baz.com", "foobarbaz")
	assert.Nil(err, "error is expected to be nil")
	assert.Equal("Foo", user.FirstName, "first name is incorrect")
}

func TestLoginUserInvalidLoginCredentials(t *testing.T) {
	assert := assert.New(t)
	respErr := errors.NewNotFoundError(fmt.Errorf("Could not find user with username and password"))
	responder := httpmock.NewJsonResponderOrPanic(respErr.Status, respErr)
	httpmock.RegisterResponder("POST", "http://localhost:8081/users/login", responder)
	repo := usersRepository{}
	user, restErr := repo.LoginUser("foo.bar@baz.com", "foobarbaz")
	assert.EqualValues("not_found", restErr.Error, "should be a not found error")
	assert.Nil(user, "user should be nil")
}

func TestLoginUserInternalServerError(t *testing.T) {
	assert := assert.New(t)
	respErr := errors.NewInternalServerError(fmt.Errorf("Could not read user from database"))
	responder := httpmock.NewJsonResponderOrPanic(respErr.Status, respErr)
	httpmock.RegisterResponder("POST", "http://localhost:8081/users/login", responder)
	repo := usersRepository{}
	user, restErr := repo.LoginUser("foo.bar@baz.com", "foobarbaz")
	assert.EqualValues("internal_server_error", restErr.Error, "should be an internal server error")
	assert.Nil(user, "user should be nil")
}

func TestLoginUserConnectionFailure(t *testing.T) {
	assert := assert.New(t)
	httpmock.RegisterResponder("POST", "http://localhost:8081/users/login", httpmock.ConnectionFailure)
	repo := usersRepository{}
	user, restErr := repo.LoginUser("foo.bar@baz.com", "foobarbaz")
	callCount := httpmock.GetCallCountInfo()
	assert.Greater(callCount["POST http://localhost:8081/users/login"], 1, "should have been retried few times but was called only once")
	assert.EqualValues("internal_server_error", restErr.Error, "should be an internal server error")
	assert.Contains(restErr.Message, "Encountered an error making downstream api call", "should be an internal server error")
	assert.Nil(user, "user should be nil")
}
