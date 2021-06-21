package scripts

import (
	"fmt"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/Squwid/bytegolf/db"
	"github.com/Squwid/bytegolf/holes"
	"github.com/Squwid/bytegolf/models"
	"github.com/Squwid/go-randomizer"
	"github.com/mitchellh/mapstructure"
)

var possibleHoles holes.Holes
var possibleBGIDs = []string{
	"2b11a2d8-029e-438d-8683-5dec6155ce53", // Mine
	"2b8eb3c40e43",
	"80851c0dcb90",
	"9dde58336d20",
	"92e081337520",
	"f031f71323e3",
	"778f511bab0b",
	"039f073a7ee9",
}

func randomDBHole() (holes.Hole, error) {
	if possibleHoles == nil {
		hs, err := allHoles()
		if err != nil {
			return holes.Hole{}, err
		}
		fmt.Println("Got", len(hs), "from DB")
		possibleHoles = hs
	}

	return possibleHoles[randomizer.Number(0, len(possibleHoles)-1)], nil
}

func randomDBBGID() string {
	return possibleBGIDs[randomizer.Number(0, len(possibleBGIDs)-1)]
}

func allHoles() (holes.Holes, error) {
	query := db.HoleCollection().OrderBy("CreatedAt", firestore.Desc).Where("Active", "==", true).Limit(100)
	docs, err := db.Query(models.NewQuery(query, nil))
	if err != nil {
		return nil, err
	}

	var hs holes.Holes
	if err := mapstructure.Decode(docs, &hs); err != nil {
		return nil, err
	}

	fmt.Printf("Got %v holes\n", len(hs))
	return hs, nil
}

func TestPopulateSubmissions(t *testing.T) {
	const subAmount = 500
	for i := 0; i < subAmount; i++ {
		h, err := randomDBHole()
		if err != nil {
			panic(err)
		}

		fmt.Printf("[%v] %v\n", i, h.ID)
	}
}
