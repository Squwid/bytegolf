package firestore

import (
	"fmt"

	"github.com/Squwid/bytegolf/globals"
)

func ProfileCollection() string {
	return prefix("Profile")
}

func prefix(collection string) string {
	return fmt.Sprintf("bg_%s_%s", globals.ENV, collection)
}
