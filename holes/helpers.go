package holes

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/Squwid/bytegolf/db"
	"github.com/Squwid/bytegolf/models"
)

// Get database hole,
/*
func getDBHole(id string) (*models.HoleDB, error) {
	doc, err := fs.Client.Collection(HolesCollection).Doc(id).Get(context.Background())
	if status.Code(err) == codes.NotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Parse the database hole
	var hole models.HoleDB
	if err := doc.DataTo(&hole); err != nil {
		return nil, err
	}
	return &hole, nil
}
*/

func listHoles() (models.Holes, error) {
	ctx := context.Background()
	docs, err := db.HoleCollection().OrderBy("CreatedAt", firestore.Desc).Where("Active", "==", true).
		Limit(10).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var holes = models.Holes{}
	for _, doc := range docs {
		var hole models.Hole
		if err := doc.DataTo(&hole); err != nil {
			return nil, err
		}
		holes = append(holes, hole)
	}

	return holes, nil
}
