package compiler

import (
	"time"

	"github.com/Squwid/bytegolf/question"
	"github.com/Squwid/bytegolf/secrets"
)

const compileURI = "https://api.jdoodle.com/v1/execute"

var jdoodleClient *secrets.Client

func init() {
	jdoodleClient = secrets.Must(secrets.GetClient("JDOODLE")).(*secrets.Client)
}

// Execute ...
type Execute struct {
	Script       string `json:"script"`
	Language     string `json:"language"`
	VersionIndex string `json:"versionIndex"`
	HoleID       string `json:"holeId,omitempty"`
}

// FullSubmission is the full submission that gets stored in the database, this should
// never be returned the user because it has all of the test cases with inputs and outputs
type FullSubmission struct {
	UUID          string    `json:"uuid"`           // gets set in RunTests
	BGID          string    `json:"bgid"`           // gets set in request
	Correct       bool      `json:"correct"`        // gets set in RunTests
	HoleID        string    `json:"holeId"`         // gets set in RunTests
	Length        int       `json:"length"`         // gets set in RunTests
	SubmittedTime time.Time `json:"submitted_time"` // gets set in RunTests

	Exe         Execute      `json:"execute"`      // gets set in RunTests
	TestOutputs []TestOutput `json:"test_outputs"` // gets set in RunTests
}

// TestOutput is the output of the tests that were run,
type TestOutput struct {
	question.Test
	ExecuteResponse
	Correct bool
}

// ExecuteResponse that gets sent to the user after each request is made, includes whether it is correct, the length
// and if it is the users best score so far
type ExecuteResponse struct {
	Output     string `json:"output"`
	StatusCode int    `json:"statusCode"`
	Memory     string `json:"memory"`
	CPUTime    string `json:"cpuTime"`
}

// ShortSubmission is the submission that gets returned to the user after
// they submit a full submission to hide tests
type ShortSubmission struct {
	ID            string    `json:"id"`
	Correct       bool      `json:"correct"`
	Language      string    `json:"language"`
	Score         int       `json:"score"`
	SubmittedTime time.Time `json:"submitted_time"`
}

// TransformToShort transforms a full submission to a short submission
func (fs FullSubmission) TransformToShort() ShortSubmission {
	return ShortSubmission{
		ID:            fs.UUID,
		Correct:       fs.Correct,
		Language:      fs.Exe.Language,
		Score:         fs.Length,
		SubmittedTime: fs.SubmittedTime,
	}
}
