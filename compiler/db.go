package compiler

import (
	"time"

	"github.com/Squwid/bytegolf/db"
	"github.com/Squwid/bytegolf/models"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

// SubmissionDB is complete submission object that gets stored in the database
type SubmissionDB struct {
	ID       string
	Script   string
	Language string
	Version  string
	Correct  bool
	BGID     string
	HoleID   string
	Length   int64

	// Timestamps
	SubmittedTime time.Time

	// Tests holds data for each individual test and whether it was correct or not
	Tests     SubmissionTests
	TestCount int
}

// SubmissionTest is the test output that gets stored in the Tests field under SubmissionDB
type SubmissionTest struct {
	Correct bool   // Test correct or not
	Hidden  bool   // Whether to show the output on the frontend or not
	Output  string `json:"Output,omitempty"` // User test output
}

// ShortSubmission is the short submission for the frontend to list all submissions,
// correct or incorrect without the entire script
type ShortSubmission struct {
	ID            string
	Language      string
	Version       string
	BGID          string
	HoleID        string
	Length        int64
	SubmittedTime time.Time
	Correct       bool
	HoleName      string
}

// FullSubmission extends the ShortSubmission with the Script and TestCases
type FullSubmission struct {
	ShortSubmission

	Script string
	Tests  map[string]SubmissionTest
}

type SubmissionTests map[string]SubmissionTest

// HideHidden hides the hidden test cases after getting the entire object from the database
func (st SubmissionTests) HideHidden() SubmissionTests {
	for k := range st {
		if st[k].Hidden {
			st[k] = SubmissionTest{
				Hidden:  true,
				Correct: st[k].Correct,
			}
		}
	}
	return st
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

		SubmittedTime: time.Now().UTC(),

		// TODO: Add a better length counter specifically for bytes
		Length:    int64(len(script)),
		Tests:     make(SubmissionTests),
		TestCount: 0,
	}
}

// ShortSub turns a SubmissionDB to a ShortSubmission. If holename is set to true, the database is hit
// and the hole name is grabbed, otherwise the error check is not needed and the HoleName will be set to ""
func (sub *SubmissionDB) ShortSub(holename bool) (*ShortSubmission, error) {
	var name = ""
	if holename {
		getter := models.NewGet(db.HoleCollection().Doc(sub.HoleID), nil)
		hole, err := db.Get(getter)
		if err != nil {
			return nil, err
		}
		name = hole["Name"].(string)
	}

	return &ShortSubmission{
		ID:            sub.ID,
		Language:      sub.Language,
		Version:       sub.Version,
		BGID:          sub.BGID,
		HoleID:        sub.HoleID,
		Length:        sub.Length,
		SubmittedTime: sub.SubmittedTime,
		Correct:       sub.Correct,
		HoleName:      name,
	}, nil
}

func (sub *SubmissionDB) FullSub() (*FullSubmission, error) {
	ss, err := sub.ShortSub(true)
	if err != nil {
		return nil, err
	}

	return &FullSubmission{
		ShortSubmission: *ss,
		Script:          sub.Script,
		Tests:           sub.Tests.HideHidden(),
	}, nil
}

func (sub *SubmissionDB) AddTest(testID, output string, correct, hidden bool) {
	sub.TestCount++

	sub.Tests[testID] = SubmissionTest{
		Correct: correct,
		Hidden:  hidden,
		Output:  output,
	}

	if !correct {
		sub.Correct = false
	}
}

func (sub SubmissionDB) Entry() Entry {
	return Entry{
		ID:       sub.ID,
		Language: sub.Language,
		Version:  sub.Version,
		Length:   sub.Length,
		HoleID:   sub.HoleID,
		BGID:     sub.BGID,
	}
}

/* Store interface functions */
func (sub SubmissionDB) Collection() *firestore.CollectionRef { return db.SubmissionsCollection() }
func (sub SubmissionDB) DocID() string                        { return sub.ID }
func (sub SubmissionDB) Data() interface{}                    { return sub }
