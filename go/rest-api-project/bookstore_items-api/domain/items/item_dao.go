package items

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_items-api/clients/elasticsearch"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_items-api/domain/queries"
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

func (i *Item) Get() *rest_errors.RestErr {
	itemId := i.Id // elasticsearch marshal and unmarshal does not repopulate the id in the item

	res, err := elasticsearch.Client.Get(esItemsIndex, i.Id)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return rest_errors.NewNotFoundError(fmt.Errorf("could not find document by id - %s", i.Id))
		}
		return rest_errors.NewInternalServerError(fmt.Errorf("error fetching document from elasticsearch"))
	}

	bytes, err := res.Source.MarshalJSON()
	if err != nil {
		return rest_errors.NewInternalServerError(fmt.Errorf("could not serialize the item retrieved from elasticsearch by id %s", i.Id))
	}

	if err := json.Unmarshal(bytes, i); err != nil {
		return rest_errors.NewInternalServerError(fmt.Errorf("document with id %s received from elasticsearch does not match the structure of item", i.Id))
	}

	i.Id = itemId

	return nil
}

func (i *Item) Update() *rest_errors.RestErr {
	doc, err := i.Doc()
	if err != nil {
		return err
	}
	res, esUpdateErr := elasticsearch.Client.Update(esItemsIndex, i.Id, doc)
	if esUpdateErr != nil {
		return rest_errors.NewInternalServerError(fmt.Errorf("backend update failed"))
	}
	if !strings.Contains(res.Result, "updated") && !strings.Contains(res.Result, "noop") {
		return rest_errors.NewInternalServerError(fmt.Errorf("backend update failed"))
	}
	return nil
}

func (i *Item) Delete() *rest_errors.RestErr {
	getErr := i.Get()
	if getErr != nil {
		return getErr
	}
	res, esDeleteErr := elasticsearch.Client.Delete(esItemsIndex, i.Id)
	if esDeleteErr != nil {
		return rest_errors.NewInternalServerError(fmt.Errorf("delete failed in backend"))
	}
	if !strings.Contains(res.Result, "deleted") {
		return rest_errors.NewInternalServerError(fmt.Errorf("delete failed in backend"))
	}
	return nil
}

func Search(q queries.EsQuery) ([]Item, *rest_errors.RestErr) {
	res, err := elasticsearch.Client.Search(esItemsIndex, q.Build())
	if err != nil {
		return nil, rest_errors.NewInternalServerError(fmt.Errorf("backend search failed"))
	}
	items := make([]Item, res.TotalHits())
	for index, h := range res.Hits.Hits {
		var i Item
		bytes, marshalErr := h.Source.MarshalJSON()
		if marshalErr != nil {
			return items, rest_errors.NewInternalServerError(fmt.Errorf("could not serialize the item retrieved from elasticsearch - %v", marshalErr))
		}
		if unmarshalErr := json.Unmarshal(bytes, &i); err != nil {
			return items, rest_errors.NewInternalServerError(fmt.Errorf("document with received from elasticsearch does not match the structure of item - %v", unmarshalErr))
		}
		i.Id = h.Id // Repopulate the id from the hit because the id field is not returned in the search result source document
		items[index] = i
	}
	return items, nil
}

func (i *Item) Doc() (map[string]interface{}, *rest_errors.RestErr) {
	var doc map[string]interface{}
	bytes, err := json.Marshal(i)
	if err != nil {
		return doc, rest_errors.NewInternalServerError(fmt.Errorf("could not get json for values of item - %v", err))
	}
	if err := json.Unmarshal(bytes, &doc); err != nil {
		return doc, rest_errors.NewInternalServerError(fmt.Errorf("could not get map for values of item - %v", err))
	}
	return doc, nil
}
