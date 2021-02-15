package access_token

import (
	"fmt"
	"strings"
	"time"

	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/utils/errors"
)

const (
	expirationTime = 24
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	ClientId    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

func GetNewAccessToken() AccessToken {
	return AccessToken{
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at AccessToken) IsExpired() bool {
	now := time.Now().UTC()
	expirationTime := time.Unix(at.Expires, 0)
	return now.After(expirationTime)
}

func (at AccessToken) Validate() *errors.RestErr {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if len(at.AccessToken) == 0 {
		return errors.NewBadRequestError(fmt.Errorf("Access token id provided by user cant be empty"))
	}
	if at.UserId <= 0 {
		return errors.NewBadRequestError(fmt.Errorf("Access token user id cant be empty"))
	}
	if at.ClientId <= 0 {
		return errors.NewBadRequestError(fmt.Errorf("Access token client id cant be empty"))
	}
	if at.Expires <= 0 {
		return errors.NewBadRequestError(fmt.Errorf("Access token expiration time cant be empty"))
	}
	return nil
}
