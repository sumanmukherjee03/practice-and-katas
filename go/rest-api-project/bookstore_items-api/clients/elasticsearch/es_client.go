package elasticsearch

import (
	"context"
	"fmt"
	"time"

	"github.com/olivere/elastic"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_items-api/logger"
)

var (
	log                      = logger.GetLogger()
	Client esClientInterface = &esClient{} // Expose a var of type interface and not a concrete struct so that it's easier to mock
)

type esClientInterface interface {
	setClient(c *elastic.Client)
	Index(string, interface{}) (*elastic.IndexResponse, error)
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
		// elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		// elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
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
		logger.Error(fmt.Sprintf("encountered error when trying to index document in elasticsearch for index %s", index), err)
		return nil, err
	}
	return res, nil
}
