package items

import (
	"fmt"

	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_items-api/clients/elasticsearch"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_utils-go/rest_errors"
)

const (
	esItemsIndex = "items"
)

func init() {
	elasticsearch.Init()
}

func (i *Item) Save() *rest_errors.RestErr {
	res, err := elasticsearch.Client.Index(esItemsIndex, i)
	if err != nil {
		return rest_errors.NewInternalServerError(fmt.Errorf("backend persistence error"))
	}
	i.Id = res.Id
	return nil
}
