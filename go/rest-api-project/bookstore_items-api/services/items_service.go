package services

import (
	"net/http"

	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_items-api/domain/items"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_items-api/domain/queries"
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
	Search(queries.EsQuery) ([]items.Item, *rest_errors.RestErr)
	Update(bool, items.Item) (*items.Item, *rest_errors.RestErr)
}

type itemsService struct {
}

func (s *itemsService) Create(item items.Item) (*items.Item, *rest_errors.RestErr) {
	err := item.Save()
	if err != nil {
		return nil, &rest_errors.RestErr{Message: "Not implemented error", Status: http.StatusNotImplemented, Error: "not_implemented"}
	}
	return &item, nil
}

func (s *itemsService) Get(itemId string) (*items.Item, *rest_errors.RestErr) {
	item := items.Item{Id: itemId}
	if err := item.Get(); err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *itemsService) Search(q queries.EsQuery) ([]items.Item, *rest_errors.RestErr) {
	return items.Search(q)
}

func (s *itemsService) Update(isPartial bool, item items.Item) (*items.Item, *rest_errors.RestErr) {
	var err *rest_errors.RestErr
	currentItem, err := s.Get(item.Id)
	if err != nil {
		return nil, err
	}
	if isPartial {
		if len(item.Title) > 0 {
			currentItem.Title = item.Title
		}
		if item.Description != (items.Description{}) {
			if len(item.Description.PlainText) > 0 {
				currentItem.Description.PlainText = item.Description.PlainText
			}
			if len(item.Description.Html) > 0 {
				currentItem.Description.Html = item.Description.Html
			}
		}
		if len(item.Pictures) > 0 {
			currentItem.Pictures = append(currentItem.Pictures, item.Pictures...)
		}
		if item.Price > 0 {
			currentItem.Price = item.Price
		}
		if item.AvailableQuantity > 0 {
			currentItem.AvailableQuantity = item.AvailableQuantity
		}
		if item.SoldQuantity > 0 {
			currentItem.SoldQuantity = item.SoldQuantity
		}
		if len(item.Video) > 0 {
			currentItem.Video = item.Video
		}
	} else {
		currentItem.Title = item.Title
		currentItem.Description = item.Description
		currentItem.Pictures = item.Pictures
		currentItem.Video = item.Video
		currentItem.Price = item.Price
		currentItem.AvailableQuantity = item.AvailableQuantity
		currentItem.SoldQuantity = item.SoldQuantity
	}

	if err = currentItem.Update(); err != nil {
		return nil, err
	}

	return currentItem, nil
}
