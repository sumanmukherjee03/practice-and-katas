package db

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/clients/cassandra"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/domain/access_token"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_utils-go/rest_errors"
)

const (
	queryGetAccessToken           = "SELECT access_token, user_id, client_id, expires FROM access_tokens WHERE access_token = ?"
	queryCreateAccessToken        = "INSERT into access_tokens (access_token, user_id, client_id, expires) VALUES (?, ?, ?, ?)"
	queryUpdateAccessTokenExpires = "UPDATE access_tokens SET expires = ? WHERE access_token = ?"
)

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, *rest_errors.RestErr)
	Create(access_token.AccessToken) *rest_errors.RestErr
	UpdateExpirationTime(access_token.AccessToken) *rest_errors.RestErr
}

type dbRepository struct {
}

func NewRepository() DbRepository {
	repo := &dbRepository{}
	return repo
}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, *rest_errors.RestErr) {
	session := cassandra.GetSession()
	var res access_token.AccessToken
	if err := session.Query(queryGetAccessToken, id).Scan(
		&res.AccessToken,
		&res.UserId,
		&res.ClientId,
		&res.Expires,
	); err != nil {
		if err == gocql.ErrNotFound {
			return nil, rest_errors.NewNotFoundError(err)
		}
		return nil, rest_errors.NewInternalServerError(err)
	}

	return &res, nil
}

func (r *dbRepository) Create(at access_token.AccessToken) *rest_errors.RestErr {
	session := cassandra.GetSession()
	if err := session.Query(queryCreateAccessToken, at.AccessToken, at.UserId, at.ClientId, at.Expires).Exec(); err != nil {
		return rest_errors.NewInternalServerError(fmt.Errorf("Could not insert access_token into database - %v", err))
	}

	return nil
}

func (r *dbRepository) UpdateExpirationTime(at access_token.AccessToken) *rest_errors.RestErr {
	session := cassandra.GetSession()
	if err := session.Query(queryUpdateAccessTokenExpires, at.Expires, at.AccessToken).Exec(); err != nil {
		return rest_errors.NewInternalServerError(fmt.Errorf("Could not update access_token in database - %v", err))
	}

	return nil
}
