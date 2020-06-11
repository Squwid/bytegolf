package submissions

import (
	"time"

	"github.com/Squwid/bytegolf/question"
)

/* FULL SUBMISSION MODELS */

type FullSubmissions []FullSubmission

// FullSubmission is the full submission that gets stored in the database, this should
// never be returned the user because it has all of the test cases with inputs and outputs
type FullSubmission struct {
	Short ShortSubmission `json:"ShortSubmission"`

	// Exe         Execute      `json:"execute"`      // gets set in RunTests
	TestOutputs []TestOutput `json:"TestOutputs"` // gets set in RunTests
}

// TestOutput is the output of the tests that were run,
type TestOutput struct {
	question.Test
	// ExecuteResponse
	Correct bool
}

/* SHORT SUBMISSION THINGS */
type ShortSubmissions []ShortSubmission

// ShortSubmission is the submission that gets returned to the user after
// they submit a full submission to hide tests
type ShortSubmission struct {
	UUID          string    `json:"UUID"`
	BGID          string    `json:"BGID"`
	Correct       bool      `json:"Correct"`
	HoleID        string    `json:"HoleID"`
	Language      string    `json:"Language"`
	Length        int       `json:"Length"`
	SubmittedTime time.Time `json:"SubmittedTime"`
}

func (fss FullSubmissions) ToShortSubmissions() ShortSubmissions {
	var ss ShortSubmissions
	for _, s := range fss {
		ss = append(ss, s.Short)
	}
	return ss
}
