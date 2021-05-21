package models

import "cloud.google.com/go/firestore"

type Queryer interface {
	Query() firestore.Query
	Transform(map[string]interface{}) error
}

type Query struct {
	query     firestore.Query
	transform func(map[string]interface{}) error
}

func (q Query) Query() firestore.Query { return q.query }

func (q Query) Transform(item map[string]interface{}) error {
	if q.transform == nil {
		return nil
	}
	return q.transform(item)
}

func NewQuery(q firestore.Query, transform func(map[string]interface{}) error) Query {
	return Query{
		query:     q,
		transform: transform,
	}
}
