package queries

import (
	"github.com/olivere/elastic"
)

func (q EsQuery) Build() elastic.Query {
	query := elastic.NewBoolQuery()
	queries := make([]elastic.Query, 0)
	for _, fv := range q.Equals {
		temp := elastic.NewMatchQuery(fv.Field, fv.Value)
		queries = append(queries, temp)
	}
	query.Must(queries...)
	return query
}
