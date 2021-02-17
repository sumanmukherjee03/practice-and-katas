package rest

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/domain/users"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/utils/errors"
)

var (
	usersRestClient = rest.RequestBuilder{
		BaseURL: "https://localhost:8080",
		Timeout: 100 * time.Millisecond,
	}
)

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, *errors.RestErr)
}

type usersRepository struct {
}

func (u *usersRepository) LoginUser(email string, password string) (*users.User, *errors.RestErr) {
	req := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}
	resp := usersRestClient.Post("/users/login", req)
	if resp == nil || resp.Response == nil {
		return nil, errors.NewInternalServerError(fmt.Errorf("Downstream users api is down"))
	}
	if resp.StatusCode > 299 {
		var restErr errors.RestErr
		if err := json.Unmarshal(resp.Bytes(), &restErr); err != nil {
			return nil, errors.NewInternalServerError(fmt.Errorf("invalid error interface when trying to login user"))
		}
		return nil, &restErr
	}
	var user users.User
	if err := json.Unmarshal(resp.Bytes(), &user); err != nil {
		return nil, errors.NewInternalServerError(fmt.Errorf("error when trying to unmarshall users api response"))
	}
	return &user, nil
}

func NewRepository() RestUsersRepository {
	return &usersRepository{}
}
