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
	ID        string
	HoleID    string
	BGID      string
	CreatedAt time.Time

	// Holds all the compile inputs and outputs
	Jdoodles []Jdoodle

	// Test Inputs and Outputs
	Tests []SubmissionDBTest

	MetaData SubmissionMetaData
}

// SubmissionDBTest is the embedded struct of TestInputs and TestOutputs
type SubmissionDBTest struct {
	TestInput  TestCaseInput
	TestOutput TestCaseOutput
}

// SubmissionFrontend is what gets displayed to users
type SubmissionFrontend struct {
	ID        string
	HoleID    string
	BGID      string
	CreatedAt time.Time
	Correct   bool
	Length    int
	Language  string
}

// Frontend allows for a database object to be changed into a frontend object
func (sdb SubmissionDB) Frontend() SubmissionFrontend {
	// TODO: Somehow get github user information into this object
	return SubmissionFrontend{
		ID:        sdb.ID,
		HoleID:    sdb.HoleID,
		BGID:      sdb.BGID,
		CreatedAt: sdb.CreatedAt,
		Correct:   sdb.MetaData.Correct,
		Length:    sdb.MetaData.Length,
		Language:  sdb.MetaData.Language,
	}
}

// SubmissionTransform transforms a list of database objects to frontend objects
func SubmissionTransform(subs []SubmissionDB) []SubmissionFrontend {
	var frontends = []SubmissionFrontend{}
	for _, sub := range subs {
		frontends = append(frontends, sub.Frontend())
	}
	return frontends
}

type Jdoodle struct {
	CompileInput  CompileInput
	CompileOutput CompileOutput
}

// SubmissionMetaData is submission meta data is the data like answer being correct, and length and user
type SubmissionMetaData struct {
	Code     string // The code input of each of the tests
	Correct  bool   // If the entire submission is correct or not
	Length   int    // Length without comments
	Language string
}

// ErrLanguageNotFound is the error where a language isnt found
var ErrLanguageNotFound = errors.New("Language not found")
