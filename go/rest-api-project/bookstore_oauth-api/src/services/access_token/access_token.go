package access_token

import (
	"fmt"
	"strings"

	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/domain/access_token"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/repository/db"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/repository/rest"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_utils-go/rest_errors"
)

type Service interface {
	GetById(string) (*access_token.AccessToken, *rest_errors.RestErr)
	Create(access_token.AccessTokenRequest) (*access_token.AccessToken, *rest_errors.RestErr)
	UpdateExpirationTime(access_token.AccessToken) *rest_errors.RestErr
}

type service struct {
	dbRepository        db.DbRepository
	restUsersRepository rest.RestUsersRepository
}

// Here we are returning the interface and not the struct type
// The interface is exposed but the struct implementing the interface is not
func NewService(r db.DbRepository, ur rest.RestUsersRepository) Service {
	return &service{
		dbRepository:        r,
		restUsersRepository: ur,
	}
}

func (s *service) GetById(atId string) (*access_token.AccessToken, *rest_errors.RestErr) {
	atId = strings.TrimSpace(atId)
	if len(atId) == 0 {
		return nil, rest_errors.NewBadRequestError(fmt.Errorf("Access token id provided by user is invalid"))
	}
	return s.dbRepository.GetById(atId)
}

func (s *service) Create(atr access_token.AccessTokenRequest) (*access_token.AccessToken, *rest_errors.RestErr) {
	if err := atr.Validate(); err != nil {
		return nil, err
	}

	user, loginRestErr := s.restUsersRepository.LoginUser(atr.Username, atr.Password)
	if loginRestErr != nil {
		return nil, loginRestErr
	}

	at := access_token.GetNewAccessToken(user.Id)
	if atCreateErr := s.dbRepository.Create(at); atCreateErr != nil {
		return nil, atCreateErr
	}

	return &at, nil
}

func (s *service) UpdateExpirationTime(at access_token.AccessToken) *rest_errors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.dbRepository.UpdateExpirationTime(at)
}
