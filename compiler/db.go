package compiler

import (
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Squwid/bytegolf/db"
	"github.com/Squwid/bytegolf/models"
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

	// Timestamps
	SubmittedTime time.Time

	// Tests holds data for each individual test and whether it was correct or not
	Tests     map[string]bool
	TestCount int
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

// FullSubmission is the full submission for the frontend including the script and all short submission data
type FullSubmission struct {
	ShortSubmission

	Script string
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
		Tests:     make(map[string]bool),
		TestCount: 0,
	}
}

func (sub *SubmissionDB) ShortSub() (*ShortSubmission, error) {
	getter := models.NewGet(db.HoleCollection().Doc(sub.HoleID), nil)
	hole, err := db.Get(getter)
	if err != nil {
		return nil, err
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
		HoleName:      hole["Name"].(string),
	}, nil
}

func (sub *SubmissionDB) FullSub() (*FullSubmission, error) {
	ss, err := sub.ShortSub()
	if err != nil {
		return nil, err
	}

	return &FullSubmission{
		ShortSubmission: *ss,
		Script:          sub.Script,
	}, nil
}

func (sub *SubmissionDB) AddTest(testID string, correct bool) {
	sub.TestCount++
	sub.Tests[testID] = correct

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
