package scripts

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Squwid/bytegolf/db"
	"github.com/Squwid/bytegolf/holes"
	"github.com/Squwid/bytegolf/models"
	"github.com/Squwid/go-randomizer"
)

var difficulties = []string{
	"IMPOSSIBLE",
	"HARD",
	"MEDIUM",
	"EASY",
	"BEGINNER",
}

func randomDifficulty() string {
	return difficulties[randomizer.Number(0, len(difficulties)-1)]
}

func randomID() string {
	return fmt.Sprintf("%s_%s", randomizer.Word(), randomizer.Word())
}

func randomParagraph(min, max int) string {
	length := randomizer.Number(min, max)
	words := randomizer.Words(length)

	return strings.Join(words, " ")
}

func randomDate() time.Time {
	begin := time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 1, 0, 0, 0, 0, 0, time.UTC)
	return randomizer.Date(begin, end)
}

func randomHole() holes.Hole {
	return holes.Hole{
		ID:            randomID(),
		Name:          strings.Title(randomParagraph(3, 10)),
		Difficulty:    randomDifficulty(),
		Question:      randomParagraph(5, 20),
		CreatedAt:     randomDate(),
		CreatedBy:     "BEN",
		LastUpdatedAt: randomDate(),
		Active:        true,
	}
}

func randomTest() holes.Test {
	return holes.Test{
		ID:          randomID(),
		Input:       fmt.Sprintf("%s %s %s", randomizer.Word(), randomizer.Word(), randomizer.Word()),
		OutputRegex: "some_regex",
		Active:      true,
		CreatedAt:   randomDate(),
	}
}

func randomTests(num int) holes.Tests {
	var tests = make(holes.Tests, num)
	for i := 0; i < num; i++ {
		tests[i] = randomTest()
	}
	return tests
}

func randomHoles(num int) holes.Holes {
	var holes = make(holes.Holes, num)
	for i := 0; i < num; i++ {
		holes[i] = randomHole()
	}
	return holes
}

// Test-prefixed just for vscode to be happy
func TestPopulateHoles(t *testing.T) {
	const holeCount = 30
	holes := randomHoles(holeCount)

	for _, hole := range holes {
		if err := db.Store(hole); err != nil {
			t.Logf("Error storing hole: %v\n", err)
			t.FailNow()
		}
	}
}

func TestPopulateTests(t *testing.T) {
	var min, max = 1, 5

	query := db.HoleCollection().OrderBy("CreatedAt", firestore.Desc).Where("Active", "==", true).Limit(1)

	hs, err := db.Query(models.NewQuery(query, nil))
	if err != nil {
		t.FailNow()
	}

	for _, obj := range hs {
		id := obj["ID"].(string)
		fmt.Println("adding tests to", id)
		testCount := randomizer.Number(min, max)

		for i := 0; i < testCount; i++ {
			test := randomTest()
			_, err := db.TestSubCollection(id).Doc(test.ID).Set(context.Background(), test)
			if err != nil {
				t.FailNow()
			}
			fmt.Println("stored test", test.ID)
		}
	}
}

func TestGetTests(t *testing.T) {
	hole := "broad_unsightly"
	tests, err := holes.GetTests(hole)
	if err != nil {
		fmt.Println("error getting tests", err)
		t.FailNow()
	}

	fmt.Println("GOT TESTS", len(tests))

	for _, test := range tests {
		fmt.Printf("%+v\n", test)
	}
}
