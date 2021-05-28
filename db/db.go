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

// TestSubCollection returns a CollectionRef of the test subcollection for a hole
func TestSubCollection(hole string) *firestore.CollectionRef {
	return HoleCollection().Doc(hole).Collection(prefix("Test"))
}

func SubmissionsCollection() *firestore.CollectionRef {
	return client.Collection(prefix("Submission"))
}

func prefix(collection string) string {
	return fmt.Sprintf("bg_%s_%s", globals.ENV, collection)
}
