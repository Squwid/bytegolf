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
