package models

import "cloud.google.com/go/firestore"

// Storer is the interface for database objects that get stored
type Storer interface {
	Collection() *firestore.CollectionRef
	DocID() string
	Data() interface{}
}
