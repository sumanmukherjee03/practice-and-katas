package queries

// EsQuery sample request would be like so : {"equals": [{"field": "status", "value": "pending"}, {"field": "seller", "value": 1}]}
type EsQuery struct {
	Equals []FieldValue `json:"equals"`
}

type FieldValue struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}
