package jdoodle

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Squwid/bytegolf/globals"
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

type validation struct {
	valid bool
	msg   string // msg exists if validation is invalid

	jdoodle globals.JdoodleLang
}

func (in UserInput) Input(stdIn string) *Input {
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

func (in UserInput) validate() validation {
	var v validation

	// Make sure that the language and version match up
	jdoodle := globals.GetLanguage(in.Language, in.Version)
	if jdoodle == nil {
		v.msg = "invalid language"
		return v
	}

	if in.Script == "" {
		v.msg = "invalid script"
		return v
	}

	v.valid = true
	v.jdoodle = *jdoodle
	return v
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
