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

type Tests []Test

type Test struct {
	ID          string `json:"ID"`
	Input       string `json:"Input"`
	OutputRegex string `json:"OutputRegex"`
	Active      bool   `json:"Active"`

	CreatedAt time.Time `json:"CreatedAt"`
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
