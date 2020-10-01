package models

import (
	"errors"
	"time"
)

// IncomingSubmission is the submission coming in from the logged in user
type IncomingSubmission struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

// CompileInput is the input of Jdoodle minus the ClientID and ClientSecret
type CompileInput struct {
	Code         string `json:"script"`
	StdIn        string `json:"stdIn"`
	Language     string `json:"language"`
	VersionIndex string `json:"versionIndex"`
}

type CompileOutput struct {
	Output     string `json:"output"`
	StatusCode int    `json:"statusCode"`
	Memory     string `json:"memory"`
	CPUTime    string `json:"cpuTime"`
}

// CompileDB stores just the compiles in the database. These dont ~need~ to be
// stored but i might need them later
type CompileDB struct {
	ID        string        `json:"ID"`
	Input     CompileInput  `json:"input"`
	Output    CompileOutput `json:"output"`
	BGID      string        `json:"BGID"`
	CreatedAt time.Time     `json:"createdAt"`
}

// SubmissionDB is the struct for submissions in the database. This is NOT passed to the user
type SubmissionDB struct {
	ID     string
	HoleID string
	BGID   string

	Jdoodles []Jdoodle

	Tests       []TestCaseInput
	TestOutputs []TestCaseOutput

	MetaData SubmissionMetaData
}

type Jdoodle struct {
	CompileInput  CompileInput
	CompileOutput CompileOutput
}

// SubmissionMetaData is submission meta data is the data like answer being correct, and length and user
type SubmissionMetaData struct {
	Code    string // The code input of each of the tests
	Correct bool   // If the entire submission is correct or not
	Length  int    // Length without comments
}

// ErrLanguageNotFound is the error where a language isnt found
var ErrLanguageNotFound = errors.New("Language not found")
