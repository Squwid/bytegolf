package bg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Squwid/bytegolf/holes"
	"github.com/Squwid/bytegolf/models"
)

// UserInput is structure of the object that the user sends the backend
// which is parsed and transformed into an Input object and sent to Jdoodle.
type UserInput struct {
	Script   string `json:"script"`
	Language string `json:"language"`
	Version  string `json:"version"`
}

// Input is the input to the compiler.
type Input struct {
	Language

	Script string `json:"script"`
	Count  int    `json:"count"`
	StdIn  string `json:"stdin,omitempty"`

	response chan models.RemoteCompilerOutput
}

// Output is what comes back from the Bytegolf compiler.
type Output struct {
	StdOut   string `json:"stdout"`
	StdErr   string `json:"stderr"`
	Duration int    `json:"duration_ms"`
	TimedOut bool   `json:"timed_out"`
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

	language *Language
}

func (in UserInput) validate() validation {
	var v validation

	lang := GetLanguage(in.Language, in.Version)
	if lang == nil {
		v.msg = fmt.Sprintf("invalid language: %v:%v", in.Language, in.Version)
		return v
	}

	if in.Script == "" {
		v.msg = "invalid script"
		return v
	}

	v.valid = true
	v.language = lang
	return v
}

func (in UserInput) Input(stdIn string, language Language) *Input {
	return &Input{
		StdIn:    stdIn,
		Language: language,
		Script:   in.Script,
		Count:    1,
		response: make(chan models.RemoteCompilerOutput, 1),
	}
}

/* Interface things for the compiler interface */

func (in Input) Request() (*http.Request, error) {
	bs, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	return http.NewRequest("POST", "http://compiler.byte.golf:8080/compile",
		bytes.NewReader(bs))
}
func (in Input) Client() *http.Client {
	return &http.Client{
		Timeout: 60 * time.Second,
	}
}
func (in Input) ResponseChan() chan models.RemoteCompilerOutput {
	return in.response
}
