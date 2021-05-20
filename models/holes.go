package models

import (
	"strings"
	"time"
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

// HoleTitle sets the hole title to an id using string lower
func HoleTitle(str string) string {
	return strings.ToLower(strings.ReplaceAll(str, " ", "-"))
}
