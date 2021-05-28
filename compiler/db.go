package compiler

import (
	"cloud.google.com/go/firestore"
	"github.com/Squwid/bytegolf/db"
	"github.com/google/uuid"
)

type SubmissionDB struct {
	ID       string
	Script   string
	Language string
	Version  string
	Correct  bool
	BGID     string
	HoleID   string
	Length   int64

	// Tests holds data for each individual test and whether it was correct or not
	Tests     map[string]bool
	TestCount int
}

func NewSubmissionDB(holeID, bgID, script, language, version string) *SubmissionDB {
	return &SubmissionDB{
		ID:       uuid.New().String(),
		Script:   script,
		Language: language,
		Version:  version,
		Correct:  true,
		BGID:     bgID,
		HoleID:   holeID,

		// TODO: Add a better length counter specifically for bytes
		Length:    int64(len(script)),
		Tests:     make(map[string]bool),
		TestCount: 0,
	}
}

func (sub *SubmissionDB) AddTest(testID string, correct bool) {
	sub.TestCount++
	sub.Tests[testID] = correct

	if !correct {
		sub.Correct = false
	}
}

/* Store interface functions */
func (sub SubmissionDB) Collection() *firestore.CollectionRef { return db.SubmissionsCollection() }
func (sub SubmissionDB) DocID() string                        { return sub.ID }
func (sub SubmissionDB) Data() interface{}                    { return sub }
