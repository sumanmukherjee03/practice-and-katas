package elasticsearch

import (
	"context"
	"fmt"
	"time"

	"github.com/olivere/elastic"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_utils-go/logger"
)

var (
	log                      = logger.GetLogger()
	Client esClientInterface = &esClient{} // Expose a var of type interface and not a concrete struct so that it's easier to mock
)

type esClientInterface interface {
	setClient(c *elastic.Client)
	Index(string, interface{}) (*elastic.IndexResponse, error)
	Get(string, string) (*elastic.GetResult, error)
	Search(string, elastic.Query) (*elastic.SearchResult, error)
}

type esClient struct {
	client *elastic.Client
}

// We have Init and not the auto init func because we dont want to connect to an actual elasticsearch
// cluster in test mode. This Init allows us to call itself on demand.
// Note : Disable sniffing when running elasticsearch in docker container - https://github.com/olivere/elastic/issues/312
func Init() {
	c, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetGzip(true),
		elastic.SetSniff(false),
		elastic.SetErrorLog(log),
		elastic.SetInfoLog(log),
	)

	if err != nil {
		panic(err)
	}

	Client.setClient(c)
}

func (c *esClient) setClient(ec *elastic.Client) {
	c.client = ec
}

func (c *esClient) Index(index string, doc interface{}) (*elastic.IndexResponse, error) {
	ctx := context.Background()
	res, err := c.client.Index().
		Index(index).
		BodyJson(doc).
		Do(ctx)

	if err != nil {
		log.Error(fmt.Sprintf("encountered error when trying to index document in elasticsearch for index %s", index), err)
		return nil, err
	}
	return res, nil
}

func (c *esClient) Get(index string, id string) (*elastic.GetResult, error) {
	ctx := context.Background()
	res, err := c.client.Get().
		Index(index).
		Id(id).
		Do(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("encountered error when trying to get a document by id from the index in elasticsearch - id : %s, index : %s", id, index), err)
		return nil, err
	}
	return res, nil
}

func (c *esClient) Search(index string, query elastic.Query) (*elastic.SearchResult, error) {
	ctx := context.Background()
	res, err := c.client.Search(index).
		Query(query).
		Do(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("encountered error when trying to search the index in elasticsearch - index : %s, query : %s", index, query), err)
		return nil, err
	}
	return res, nil
}
