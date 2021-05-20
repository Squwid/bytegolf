package db

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/Squwid/bytegolf/globals"
)

func ProfileCollection() *firestore.CollectionRef {
	return Client.Collection(prefix("Profile"))
}

func HoleCollection() *firestore.CollectionRef {
	return Client.Collection(prefix("Hole"))
}

func prefix(collection string) string {
	return fmt.Sprintf("bg_%s_%s", globals.ENV, collection)
}
