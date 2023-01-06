package models

import "cloud.google.com/go/firestore"

type Getter interface {
	Doc() *firestore.DocumentRef
	Transform(map[string]interface{}) error
}

type Get struct {
	doc       *firestore.DocumentRef
	transform func(map[string]interface{}) error
}

func (g Get) Doc() *firestore.DocumentRef { return g.doc }

func (g Get) Transform(item map[string]interface{}) error {
	if g.transform == nil {
		return nil
	}
	return g.transform(item)
}

func NewGet(doc *firestore.DocumentRef, transform func(map[string]interface{}) error) Get {
	return Get{
		doc:       doc,
		transform: transform,
	}
}

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

// Storer is the interface for database objects that get stored
type Storer interface {
	Collection() *firestore.CollectionRef
	DocID() string
	Data() interface{}
}
