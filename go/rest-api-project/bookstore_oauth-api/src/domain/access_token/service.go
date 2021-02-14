package access_token

import (
	"fmt"
	"strings"

	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/utils/errors"
)

// This interface will be satisfied by the database repository
type Repository interface {
	GetById(string) (*AccessToken, *errors.RestErr)
}

type Service interface {
	GetById(string) (*AccessToken, *errors.RestErr)
}

type service struct {
	repository Repository
}

// Here we are returning the interface and not the struct type
// The interface is exposed but the struct implementing the interface is not
func NewService(repo Repository) Service {
	return &service{
		repository: repo,
	}
}

func (s *service) GetById(accessTokenId string) (*AccessToken, *errors.RestErr) {
	atId := strings.TrimSpace(accessTokenId)
	if len(atId) == 0 {
		return nil, errors.NewBadRequestError(fmt.Errorf("Access token id provided by user is invalid"))
	}
	return s.repository.GetById(accessTokenId)
}
