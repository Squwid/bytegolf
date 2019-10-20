package firestore

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
)

var projectID = os.Getenv("PROJECT_ID")

// Client holds teh firestore client that is used in this api
var Client *firestore.Client

func init() {
	ctx := context.Background()
	c, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	Client = c
}

// StoreData in a particular collection
func StoreData(collection, id string, data interface{}) error {
	ctx := context.Background()
	_, _, err := Client.Collection(collection).Add(ctx, data)
	return err
}
