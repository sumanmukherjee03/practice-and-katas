package rest

import (
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/domain/users"
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
	responder := httpmock.NewJsonResponderOrPanic(200, respUser)
	httpmock.RegisterResponder("POST", "http://localhost:8080/users/login", responder)
	repo := usersRepository{}
	res, err := repo.LoginUser("foo.bar@baz.com", "foobarbaz")
	assert.Nil(err, "error is expected to be nil")
	assert.Equal("Foo", res.FirstName, "first name is incorrect")
}

// func TestLoginUserInvalidLoginCredentials(t *testing.T) {
// }

// func TestLoginUserInvalidUserJsonResponse(t *testing.T) {
// }

// func TestLoginUserNoError(t *testing.T) {
// }
