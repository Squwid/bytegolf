package firestore

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
)

var projectID = os.Getenv("PROJECT_ID")
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
func StoreData(collection string, data interface{}) error {
	ctx := context.Background()
	_, _, err := Client.Collection(collection).Add(ctx, data)
	return err
}
