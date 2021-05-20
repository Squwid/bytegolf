package db

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/Squwid/bytegolf/globals"
)

func ProfileCollection() *firestore.CollectionRef {
	return client.Collection(prefix("Profile"))
}

func HoleCollection() *firestore.CollectionRef {
	return client.Collection(prefix("Hole"))
}

func prefix(collection string) string {
	return fmt.Sprintf("bg_%s_%s", globals.ENV, collection)
}
