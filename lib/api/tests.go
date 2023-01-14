package api

import (
	"time"

	"github.com/uptrace/bun"
)

type TestsDB []TestDB

type TestDB struct {
	bun.BaseModel `bun:"table:tests,alias:t"`

	ID          int64  `bun:"id,pk,autoincrement,notnull"`
	Name        string `bun:"name,notnull"`
	Hole        string `bun:"hole,notnull"`
	Hidden      bool   `bun:"hidden,notnull"`
	Description string `bun:"description,notnull"`
	// Benchmark is ture if its the main test that
	// is used to calculate the CPU + Memory score.
	Benchmark bool `bun:"benchmark,notnull"`
	Active    bool `bun:"active,notnull"`

	Input       string `bun:"input"`
	OutputRegex string `bun:"regex,notnull"`

	CreatedAt time.Time `bun:"created_at,notnull"`
}

type TestClient struct {
	ID          int64  `json:"ID"`
	Name        string `json:"Name"`
	Hole        string `json:"Hole"`
	Hidden      bool   `json:"Hidden"`
	Description string `json:"Description"`
	Input       string `json:"Input,omitempty"`
}

// toClient already assumes that the tests are active.
func (tdb TestDB) toClient() TestClient {
	tc := TestClient{
		ID:          tdb.ID,
		Name:        tdb.Name,
		Hole:        tdb.Hole,
		Hidden:      tdb.Hidden,
		Description: tdb.Description,
	}
	if !tdb.Hidden {
		tc.Input = tdb.Input
	}
	return tc
}
