package holes

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	fs "github.com/Squwid/bytegolf/firestore"
	"github.com/Squwid/bytegolf/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HolesCollection - db collection
const HolesCollection = "holes"

// Get database hole,
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

// GetDBHole gets a DB hole, which contains tests and things
func GetDBHole(id string) (*models.HoleDB, error) {
	return getDBHole(id)
}

func storeDBHole(hole *models.HoleDB) error {
	hole.LastUpdatedAt = time.Now()
	_, err := fs.Client.Collection(HolesCollection).Doc(hole.Hole.ID).Set(context.Background(), *hole)
	return err
}

// GetHole gets the hole that doesnt include the database hole
func GetHole(id string, showInactive bool) (*models.Hole, error) {
	dbHole, err := getDBHole(id)
	if err != nil {
		return nil, err
	}

	if dbHole == nil {
		return nil, nil
	}

	if !dbHole.Active && !showInactive {
		return nil, nil
	}

	return &dbHole.Hole, nil
}

// GetAllHoles gets all the holes, active or not
func getAllDBHoles(onlyActive bool) ([]models.HoleDB, error) {
	var iter *firestore.DocumentIterator

	col := fs.Client.Collection(HolesCollection)
	if onlyActive {
		iter = col.Where("Active", "==", true).Documents(context.Background())
	} else {
		iter = col.Documents(context.Background())
	}

	// Get all of the documents
	docs, err := iter.GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return []models.HoleDB{}, nil
	}

	var holes = []models.HoleDB{}
	for _, doc := range docs {
		// Parse each hole, and append to holes
		var hole models.HoleDB
		if err := doc.DataTo(&hole); err != nil {
			return nil, err
		}

		holes = append(holes, hole)
	}

	return holes, nil
}

// GetHoles get the holes that dont include database objects
func GetHoles(onlyActive bool) ([]models.Hole, error) {
	dbHoles, err := getAllDBHoles(onlyActive)
	if err != nil {
		return nil, err
	}

	var holes = []models.Hole{}
	for _, dbHole := range dbHoles {
		holes = append(holes, dbHole.Hole)
	}

	return holes, nil
}

// ErrHoleDoesntExist gets returned when trying to modify a hole that doesnt exist
// var ErrHoleDoesntExist = errors.New("Hole doesnt exist")

// func addTestToHole(holeID string) error {
// 	// Check if hole exists

// 	hole, err := getDBHole(holeID)
// 	if err != nil {
// 		return err
// 	}
// 	if hole == nil {
// 		return ErrHoleDoesntExist
// 	}

// }
