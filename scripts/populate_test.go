package scripts

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Squwid/bytegolf/models"
	"github.com/Squwid/go-randomwords"
)

var difficulties = []string{
	"IMPOSSIBLE",
	"HARD",
	"MEDIUM",
	"EASY",
	"BEGINNER",
}

func randomDifficulty() string {
	return difficulties[randomwords.Number(0, len(difficulties)-1)]
}

func randomID() string {
	return fmt.Sprintf("%s_%s", randomwords.Word(), randomwords.Word())
}

func randomParagraph(min, max int) string {
	length := randomwords.Number(min, max)
	words := randomwords.Words(length)

	return strings.Join(words, " ")
}

func randomDate() time.Time {
	begin := time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 1, 0, 0, 0, 0, 0, time.UTC)
	return randomwords.Date(begin, end)
}

func randomHole() models.Hole {
	return models.Hole{
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

func randomHoles(num int) models.Holes {
	var holes = make(models.Holes, num)
	for i := 0; i < num; i++ {
		holes[i] = randomHole()
	}
	return holes
}

// Test-prefixed just for vscode to be happy
func TestPopulateHoles(t *testing.T) {
	const holeCount = 10
	holes := randomHoles(holeCount)

	bs, err := json.Marshal(holes)
	if err != nil {
		t.FailNow()
	}

	fmt.Println(string(bs))
}
