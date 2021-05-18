package firestore

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
)

var projectID string

// Client holds teh firestore client that is used in this api
var Client *firestore.Client

func init() {
	fmt.Println("creds: " + os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))

	projectID = os.Getenv("GCP_PROJECT_ID")

	c, err := firestore.NewClient(context.Background(), projectID)
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
