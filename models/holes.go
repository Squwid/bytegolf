package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Hole is frontend hole structure
type Hole struct {
	// ID has to be no spaces, alphanumeric only
	ID         string
	Name       string
	Difficulty string
	Question   string
}

// HoleDB inherits Hole with extra database fields. Dont export this to the user
type HoleDB struct {
	Hole Hole

	CreatedAt     time.Time
	CreatedBy     string
	LastUpdatedAt time.Time
	Active        bool
}

// NewHoleDB creates a new hole object for the database. If id is provided, it will be used.
// otherwise a new id will be generated
func NewHoleDB(id, name, difficulty, question string) *HoleDB {
	// if id is not provided, generate a new one
	if id == "" {
		id = uuid.New().String()
	} else {
		// strip spaces and move to lowercase
		id = strings.ToLower(id)
		id = strings.ReplaceAll(id, " ", "-")
	}
	return &HoleDB{
		Hole: Hole{
			ID:         id,
			Name:       name,
			Difficulty: difficulty,
			Question:   question,
		},
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		Active:        true,
	}
}

// HoleTitle sets the hole title to an id using string lower
func HoleTitle(str string) string {
	return strings.ToLower(strings.ReplaceAll(str, " ", "-"))
}
