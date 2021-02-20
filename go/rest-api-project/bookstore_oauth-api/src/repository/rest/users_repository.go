package rest

import (
	"fmt"
	"time"

	"github.com/go-resty/resty"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/domain/users"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/utils/errors"
)

var (
	usersRestClient = resty.New()
)

func init() {
	usersRestClient.SetTimeout(100*time.Millisecond).
		SetHeader("Accept", "application/json").
		SetHeader("User-Agent", "go-resty").
		SetHeader("X-Public", "false").
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(20 * time.Second).
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 0, fmt.Errorf("quota exceeded")
		})
}

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, *errors.RestErr)
}

type usersRepository struct {
}

func (u *usersRepository) LoginUser(email string, password string) (*users.User, *errors.RestErr) {
	var user users.User
	var restErr errors.RestErr
	req := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}
	_, err := GetRestClient().R().
		SetBody(req).
		SetResult(&user).
		SetError(&restErr).
		Post("http://localhost:8080/users/login")
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Errorf("Encountered an error making downstream api call - %v", err))
	}
	if restErr.Status > 0 {
		return nil, &restErr
	}
	return &user, nil
}

func NewRepository() RestUsersRepository {
	return &usersRepository{}
}

func GetRestClient() *resty.Client {
	return usersRestClient
}
