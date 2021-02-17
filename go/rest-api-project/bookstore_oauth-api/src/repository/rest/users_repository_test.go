package rest

import (
	"net/http"
	"os"
	"testing"

	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
)

// TestMain is the entrypoint of every test suite in a package is the TestMain function
// just like Main function is the entrypoint for a go package
func TestMain(m *testing.M) {
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromApi(t *testing.T) {
	assert := assert.New(t)
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   "POST",
		URL:          "http://localhost:8080/users/login",
		ReqBody:      `{"email":"foo@bar.com","password":"foobarbaz"}`,
		RespHTTPCode: -1,
	})
	repo := usersRepository{}
	user, err := repo.LoginUser("foo@bar.com", "foobarbaz")
	assert.Nil(user, "user should have been nil")
	assert.NotNil(err, "err should not have been nil")
	assert.EqualValues(err.Status, http.StatusInternalServerError, "err should have been internal server error")
}

// func TestLoginUserInvalidErrorInterface(t *testing.T) {
// }

// func TestLoginUserInvalidLoginCredentials(t *testing.T) {
// }

// func TestLoginUserInvalidUserJsonResponse(t *testing.T) {
// }

// func TestLoginUserNoError(t *testing.T) {
// }
