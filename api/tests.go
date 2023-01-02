package api

import (
	"time"

	"github.com/uptrace/bun"
)

type TestsDB []TestDB

type TestDB struct {
	bun.BaseModel `bun:"table:tests"`

	ID          int64  `bun:"id,pk,autoincrement,notnull"`
	Name        string `bun:"name,notnull"`
	Hole        string `bun:"hole,notnull"`
	Hidden      bool   `bun:"hidden,notnull"`
	Description string `bun:"description,notnull"`
	Active      bool   `bun:"active,notnull"`

	Input       string `bun:"input"`
	OutputRegex string `bun:"regex,notnull"`

	CreatedAt time.Time `bun:"created_at,notnull"`
}
