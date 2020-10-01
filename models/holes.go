package models

import (
	"strings"
	"time"
)

// Hole is frontend hole structure
type Hole struct {
	// ID has to be no spaces, alphanumeric only
	ID         string `json:"ID"`
	Name       string `json:"Name"`
	Difficulty string `json:"Difficulty"`
	Question   string `json:"Question"`
}

// HoleDB inherits Hole with extra database fields. Dont export this to the user
type HoleDB struct {
	Hole Hole `json:"Hole"`

	CreatedAt     time.Time `json:"CreatedAt"`
	CreatedBy     string    `json:"CreatedBy"`
	LastUpdatedAt time.Time `json:"LastUpdatedAt"`
	Active        bool      `json:"Active"`

	// Test Cases
	Tests []TestCaseInput `json:"Tests"`
}

// TestCaseInput is the struct for each of the test cases
type TestCaseInput struct {
	ID    string `json:"ID"`
	Input string `json:"Input"`
	// Solution is the test case solution, in Regex
	Solution string `json:"Solution"`
}

type TestCaseOutput struct {
	ID      string `json:"ID"` // Matches the TestCaseInputID
	Output  string `json:"Output"`
	Correct bool   `json:"Correct"`
}

// HoleTitle sets the hole title to an id using string lower
func HoleTitle(str string) string {
	return strings.ToLower(strings.ReplaceAll(str, " ", "-"))
}
