package elasticsearch

import (
	"context"
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
	Index(interface{}) (*elastic.IndexResponse, error)
}

type esClient struct {
	client *elastic.Client
}

// We have Init and not the auto init func because we dont want to connect to an actual elasticsearch
// cluster in test mode. This Init allows us to call itself on demand.
func Init() {
	c, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetGzip(true),
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

func (c *esClient) Index(data interface{}) (*elastic.IndexResponse, error) {
	return c.client.Index().Do(context.Background())
}
