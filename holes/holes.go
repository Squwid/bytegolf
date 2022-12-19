package holes

import (
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Squwid/bytegolf/db"
)

type Holes []Hole

type Hole struct {
	ID         string `json:"ID"`
	Name       string `json:"Name"`
	Difficulty string `json:"Difficulty"`
	Question   string `json:"Question"`

	CreatedAt     time.Time `json:"CreatedAt"`
	CreatedBy     string    `json:"CreatedBy"`
	LastUpdatedAt time.Time `json:"LastUpdatedAt"`
	Active        bool      `json:"Active"`
}

type ShortTests []ShortTest

// ShortTest is the frontend object for a Test structure
type ShortTest struct {
	ID          string `json:"ID"`
	Name        string `json:"Name"`                  // Test name on the frontend
	Hidden      bool   `json:"Hidden"`                // When hidden is false, output will be shown on frontend
	Input       string `json:"Input,omitempty"`       // Returns test input when test is NOT hidden
	Description string `json:"Description,omitempty"` // Optional description field
	Active      bool   `json:"Active"`
}

type Tests []Test

// Test extends ShortTest and has hidden information solely for the database
type Test struct {
	ID          string `json:"ID"`
	Name        string `json:"Name"`        // Test name on the frontend
	Hidden      bool   `json:"Hidden"`      // When hidden is false, output will be shown on frontend
	Description string `json:"Description"` // Optional description field
	Active      bool   `json:"Active"`

	Input       string `json:"Input"`
	OutputRegex string `json:"OutputRegex"`

	CreatedAt time.Time `json:"CreatedAt"`
}

func (ts Tests) ShortTests() ShortTests {
	var sts = make(ShortTests, len(ts))
	for i := 0; i < len(ts); i++ {
		sts[i] = ts[i].ShortTest()
	}
	return sts
}

// ShortTest returns what a user should see for a test case. If a test is not hidden, then
// the input is returned as well
func (t Test) ShortTest() ShortTest {
	st := ShortTest{
		ID:          t.ID,
		Name:        t.Name,
		Hidden:      t.Hidden,
		Description: t.Description,
		Active:      t.Active,
	}
	if !t.Hidden {
		st.Input = t.Input
	}
	return st
}

func transformHole(hole map[string]interface{}) error {
	delete(hole, "CreatedBy")
	return nil
}

// HoleTitle sets the hole title to an id using string lower
func HoleTitle(str string) string {
	return strings.ToLower(strings.ReplaceAll(str, " ", "-"))
}

// DB Interface stuff
func (h Hole) Collection() *firestore.CollectionRef { return db.HoleCollection() }
func (h Hole) DocID() string                        { return h.ID }
func (h Hole) Data() interface{}                    { return h }

// Sort interface stuff
func (hs Holes) Len() int           { return len(hs) }
func (hs Holes) Swap(i, j int)      { hs[i], hs[j] = hs[j], hs[i] }
func (hs Holes) Less(i, j int) bool { return hs[i].CreatedAt.Before(hs[j].CreatedAt) }
