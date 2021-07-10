package jdoodle

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Squwid/bytegolf/holes"
	"github.com/Squwid/bytegolf/models"
)

// UserInput is structure of the object that the user sends the backend
// which is parsed and transformed into an Input object and sent to Jdoodle
type UserInput struct {
	Script   string `json:"script"`
	Language string `json:"language"`
	Version  string `json:"version"`
}

// Input gets sent to Jdoodle compiler for code submission
type Input struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Script       string `json:"script"`
	StdIn        string `json:"stdIn,omitempty"`
	Language     string `json:"language"`
	VersionIndex string `json:"versionIndex"`

	response chan models.RemoteCompilerOutput
}

// Output is what comes back from the Jdoodle compiler
type Output struct {
	Output     string `json:"output"`
	StatusCode int    `json:"statusCode"`
	Memory     string `json:"memory"`
	CPUTime    string `json:"cpuTime"`
}

// compileResult is the result of the user input after all tests have been run
type compileResult struct {
	correct bool
	err     error

	output *Output
	test   *holes.Test
}

func (in UserInput) Input(stdIn string) *Input {
	// TODO: Verify that language and version can be used before calling compiler
	// TODO: StdIn needs to come from the tests somewhere
	return &Input{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Script:       in.Script,
		StdIn:        stdIn,
		VersionIndex: in.Version,
		Language:     in.Language,
		response:     make(chan models.RemoteCompilerOutput, 1),
	}
}

func (in UserInput) validate() (bool, string) {
	if in.Language == "" {
		return false, "invalid language"
	}
	if in.Script == "" {
		return false, "invalid script"
	}
	if in.Version == "" {
		return false, "invalid version"
	}
	return true, ""
}

/* Interface things for the compiler interface */

func (in Input) Request() (*http.Request, error) {
	bs, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	return http.NewRequest("POST", "https://api.jdoodle.com/v1/execute", bytes.NewReader(bs))
}
func (in Input) Client() *http.Client                           { return http.DefaultClient }
func (in Input) ResponseChan() chan models.RemoteCompilerOutput { return in.response }
