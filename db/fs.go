package db

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/Squwid/bytegolf/models"
)

var projectID string

// Client holds teh firestore client that is used in this api
var client *firestore.Client

func init() {
	fmt.Println("creds: " + os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))

	projectID = os.Getenv("GCP_PROJECT_ID")

	c, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		panic(err)
	}
	client = c
}

func Store(storer models.Storer) error {
	ctx := context.Background()
	_, err := storer.Collection().Doc(storer.DocID()).Set(ctx, storer.Data())
	return err
}

func Query(query models.Queryer) ([]map[string]interface{}, error) {
	ctx := context.Background()
	docs, err := query.Query().Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var m = make([]map[string]interface{}, len(docs))
	for i, doc := range docs {
		var item map[string]interface{}
		if err := doc.DataTo(&item); err != nil {
			return nil, err
		}

		if err := query.Transform(item); err != nil {
			return nil, err
		}

		m[i] = item
	}

	return m, nil
}

func Get(getter models.Getter) (map[string]interface{}, error) {
	ctx := context.Background()
	doc, err := getter.Doc().Get(ctx)
	if err != nil {
		return nil, err
	}

	return doc.Data(), nil
}
