package scripts

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Squwid/bytegolf/db"
	"github.com/Squwid/bytegolf/holes"
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
