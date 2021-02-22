package oauth

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-go/errors"
)

const (
	headerXPublic    = "X-Public"
	headerXClientId  = "X-Client-Id"
	headerXCallerId  = "X-Caller-Id"
	paramAccessToken = "access_token"

	maxRetries            = 10
	retryWaitTime         = 200 * time.Millisecond
	maxRetryWaitTime      = 3 * time.Second
	oAuthUserAgent        = "bookstore-oauth-go-client"
	respAcceptContentType = "application/json"
)

var (
	restClient = resty.New()
)

type oauthClient struct {
}

type oauthInterface interface {
}

type AccessToken struct {
	Id       string `json:"id"`
	UserId   int64  `json:"user_id"`
	ClientId int64  `json:"client_id"`
}

func init() {
	restClient.SetTimeout(200*time.Millisecond).
		SetHeader("Accept", respAcceptContentType).
		SetHeader("User-Agent", respAcceptContentType).
		SetHeader("X-Public", "false").
		SetRetryCount(maxRetries).
		SetRetryWaitTime(retryWaitTime).
		SetRetryMaxWaitTime(maxRetryWaitTime)
}

func IsPublic(req *http.Request) bool {
	// In case of a pointer it's always a good idea to validate if the pointer is nil or not
	if req == nil {
		return true
	}
	return req.Header.Get(headerXPublic) == "true"
}

func GetCallerId(req *http.Request) int64 {
	if req == nil {
		return 0
	}
	callerId, err := strconv.ParseInt(req.Header.Get(headerXCallerId), 10, 64)
	if err != nil {
		return 0
	}
	return callerId
}

func GetClientId(req *http.Request) int64 {
	if req == nil {
		return 0
	}
	clientId, err := strconv.ParseInt(req.Header.Get(headerXClientId), 10, 64)
	if err != nil {
		return 0
	}
	return clientId
}

func Authenticate(req *http.Request) *errors.RestErr {
	if req == nil {
		return nil
	}
	cleanRequest(req)

	token := strings.TrimSpace(req.URL.Query().Get(paramAccessToken))
	if len(token) == 0 {
		return errors.NewBadRequestError(fmt.Errorf("Access token is not present in request"))
		// return nil
	}

	at, err := getAccessToken(token)
	if err != nil {
		return err
	}
	req.Header.Add(headerXCallerId, fmt.Sprintf("%v", at.UserId))
	req.Header.Add(headerXClientId, fmt.Sprintf("%v", at.ClientId))

	return nil
}

func cleanRequest(req *http.Request) {
	if req == nil {
		return
	}
	req.Header.Del(headerXClientId)
	req.Header.Del(headerXCallerId)
}

func getAccessToken(token string) (*AccessToken, *errors.RestErr) {
	var at AccessToken
	var restErr errors.RestErr
	_, err := getRestClient().R().
		SetPathParams(map[string]string{
			"accessToken": token,
		}).
		SetResult(&at).
		SetError(&restErr).
		Get("http://localhost:8080/oauth/access_token/{accessToken}")
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Errorf("Encountered an error making downstream api call to get oauth access token - %v", err))
	}
	if restErr.Status > 0 {
		return nil, &restErr
	}
	return &at, nil
}

func getRestClient() *resty.Client {
	return restClient
}
