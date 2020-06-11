package compiler

import (
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
