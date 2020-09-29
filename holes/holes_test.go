package holes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	fs "github.com/Squwid/bytegolf/firestore"
	"github.com/Squwid/bytegolf/models"
	randomwords "github.com/Squwid/go-randomwords"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

// TestStoreHoleDB is testing to see how structures work with
// inheritance with firestore
func TestStoreHoleDB(t *testing.T) {
	// Test store, then get and see if same
	hole := models.Hole{
		ID:         uuid.New().String(),
		Name:       "Testing Hole 1",
		Difficulty: "hard",
		Question:   "This is a test question",
	}

	holedb := models.HoleDB{
		Hole:      hole,
		CreatedAt: time.Now(),
		Active:    true,
	}

	ctx := context.Background()
	_, err := fs.Client.Collection(HolesCollection).Doc(holedb.Hole.ID).Set(ctx, holedb)
	if err != nil {
		t.Logf("Error storing: %v\n", err)
		t.FailNow()
	}

	t.Logf("Stored %s\n", holedb.Hole.ID)
}

func TestMassStoreDB(t *testing.T) {
	var max, count = 10, 0

	for {
		if count >= max {
			break
		}
		count++

		randomwords.NewRandSource()

		name := fmt.Sprintf("%s %s", randomwords.RandomAdjective(), randomwords.RandomNoun())
		hole := models.NewHoleDB(name, name, "easy", "This is a question field here")

		if strings.Contains(hole.Hole.ID, "a") {
			hole.Active = false
		}

		if err := storeDBHole(hole); err != nil {
			t.Logf("Error storing DB Hole: %v\n", err)
			t.FailNow()
		}

		t.Logf("Stored %s in database\n", hole.Hole.ID)
		time.Sleep(1 * time.Second)
	}
}

func TestListAllDBHoles(t *testing.T) {
	var onlyActive = false

	holes, err := getAllDBHoles(onlyActive)
	if err != nil {
		t.Logf("Error getting db holes: %v\n", err)
		t.FailNow()
	}

	bs, err := json.Marshal(holes)
	if err != nil {
		t.Logf("Error marshalling db holes: %v\n", err)
		t.FailNow()
	}

	t.Logf("HOLES:\n%s\n", string(bs))
}

func TestListAllHoles(t *testing.T) {
	var onlyActive = false

	holes, err := GetHoles(onlyActive)
	if err != nil {
		t.Logf("Error getting db holes: %v\n", err)
		t.FailNow()
	}

	bs, err := json.Marshal(holes)
	if err != nil {
		t.Logf("Error marshalling db holes: %v\n", err)
		t.FailNow()
	}

	t.Logf("HOLES:\n%s\n", string(bs))
}

// TestGetHoleDB trys to get the hole from the databse
func TestGetHoleDB(t *testing.T) {
	const id = "f2a74539-9ac5-43aa-bfa8-1ea524594177"

	doc, err := fs.Client.Collection(HolesCollection).Doc(id).Get(context.Background())
	if err != nil {
		t.Logf("Error getting: %v\n", err)
		t.FailNow()
	}

	var hole models.HoleDB
	if err := mapstructure.Decode(doc.Data(), &hole); err != nil {
		t.Logf("Error: %v\n", err)
		t.FailNow()
	}

	fmt.Printf("%+v\n", hole)

	hBS, err := json.Marshal(hole)
	if err != nil {
		t.Logf("Error: %v\n", err)
		t.FailNow()
	}

	hhBS, err := json.Marshal(hole.Hole)
	if err != nil {
		t.Logf("Error: %v\n", err)
		t.FailNow()
	}

	t.Logf("Hole: %s\n", string(hBS))
	t.Logf("Hole.Hole: %s\n", string(hhBS))
}
