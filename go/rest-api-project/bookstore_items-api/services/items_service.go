package services

import (
	"net/http"

	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_items-api/domain/items"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_utils-go/rest_errors"
)

// Expose the itemsService struct which implements the itemsServiceInterface via
// a variable ItemsService which can then be used by other parts of the app
var (
	ItemsService itemsServiceInterface = &itemsService{}
)

type itemsServiceInterface interface {
	Create(items.Item) (*items.Item, *rest_errors.RestErr)
	Get(string) (*items.Item, *rest_errors.RestErr)
}

type itemsService struct {
}

func (s *itemsService) Create(item items.Item) (*items.Item, *rest_errors.RestErr) {
	return nil, &rest_errors.RestErr{Message: "Not implemented error", Status: http.StatusNotImplemented, Error: "not_implemented"}
}

func (s *itemsService) Get(itemId string) (*items.Item, *rest_errors.RestErr) {
	return nil, &rest_errors.RestErr{Message: "Not implemented error", Status: http.StatusNotImplemented, Error: "not_implemented"}
}
