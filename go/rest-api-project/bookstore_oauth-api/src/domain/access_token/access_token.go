package access_token

import (
	"fmt"
	"strings"
	"time"

	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_users-api/utils/crypto_utils"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_utils-go/rest_errors"
)

const (
	expirationTime             = 24
	grantTypePassword          = "password"
	grantTypeClientCredentials = "client_credentials"
)

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
	// Used for grant_type : password
	Username string `json:"username"`
	Password string `json:"password"`
	// Used for grant_type : client_credentials
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`

	Scope string `json:"scope"`
}

func (atr AccessTokenRequest) Validate() *rest_errors.RestErr {
	atr.GrantType = strings.TrimSpace(atr.GrantType)
	if len(atr.GrantType) == 0 {
		return rest_errors.NewBadRequestError(fmt.Errorf("Access token request grant type provided by user is empty"))
	}

	switch atr.GrantType {
	case grantTypePassword:
		atr.Username = strings.TrimSpace(atr.Username)
		if len(atr.Username) == 0 {
			return rest_errors.NewBadRequestError(fmt.Errorf("Access token request username is empty"))
		}
		atr.Password = strings.TrimSpace(atr.Password)
		if len(atr.Password) == 0 {
			return rest_errors.NewBadRequestError(fmt.Errorf("Access token request password is empty"))
		}
	case grantTypeClientCredentials:
		atr.ClientId = strings.TrimSpace(atr.ClientId)
		if len(atr.ClientId) == 0 {
			return rest_errors.NewBadRequestError(fmt.Errorf("Access token request client id is empty"))
		}
		atr.ClientSecret = strings.TrimSpace(atr.ClientSecret)
		if len(atr.ClientSecret) == 0 {
			return rest_errors.NewBadRequestError(fmt.Errorf("Access token request client secret is empty"))
		}
	default:
		return rest_errors.NewBadRequestError(fmt.Errorf("Access token request grant type is not supported"))
	}

	return nil
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	ClientId    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

func GetNewAccessToken(userId int64) AccessToken {
	expiresAt := time.Now().UTC().Add(expirationTime * time.Hour).Unix()
	token := crypto_utils.GetMd5(fmt.Sprintf("at-%d-%d-ran", userId, expiresAt))
	return AccessToken{
		AccessToken: token,
		UserId:      userId,
		ClientId:    1, // Hard coding this temporarily
		Expires:     expiresAt,
	}
}

func (at AccessToken) IsExpired() bool {
	now := time.Now().UTC()
	expirationTime := time.Unix(at.Expires, 0)
	return now.After(expirationTime)
}

func (at AccessToken) Validate() *rest_errors.RestErr {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if len(at.AccessToken) == 0 {
		return rest_errors.NewBadRequestError(fmt.Errorf("Access token id provided by user cant be empty"))
	}
	if at.UserId <= 0 {
		return rest_errors.NewBadRequestError(fmt.Errorf("Access token user id cant be empty"))
	}
	if at.ClientId <= 0 {
		return rest_errors.NewBadRequestError(fmt.Errorf("Access token client id cant be empty"))
	}
	if at.Expires <= 0 {
		return rest_errors.NewBadRequestError(fmt.Errorf("Access token expiration time cant be empty"))
	}
	return nil
}
