package runner

import (
	"os"
)

// Langs consist of every language currently available from the compiler
const (
	LangJava  = "java"
	LangC     = "c"
	LangCPP   = "cpp"
	LangCPP14 = "cpp14"
	LangPHP   = "php"
	LangPy2   = "python2"
	LangPy3   = "python3"
	LangRuby  = "ruby"
	LangGo    = "go"
	LangBash  = "bash"
	LangSwift = "swift"
	LangR     = "r"
	LangNode  = "nodejs"
	LangFS    = "fsharp"
)

// CodeSubmission is what gets submitted to the
type CodeSubmission struct {
	Script   string
	Language string
}

// CodeResponse is the response from the Code Runner API that gets a result
type CodeResponse struct {
	Output     string `json:"output"`
	StatusCode int    `json:"statusCode"`
	Memory     string `json:"memory"`
	CPUTime    string `json:"cpuTime"`
}

// Client is the client that allows the user to talk to the API
type Client struct {
	ID     string
	Secret string
}

// NewClient returns a new client using the environmental variables
// of RUNNER_ID for the ID and RUNNER_SECRET for the secret
func NewClient() *Client {
	return &Client{
		ID:     os.Getenv("RUNNER_ID"),
		Secret: os.Getenv("RUNNER_SECRET"),
	}
}

// NewClientWithCreds returns the credentials using a users credentials as
// an argument rather than envirmental variable
func NewClientWithCreds(id, secret string) *Client {
	return &Client{
		ID:     id,
		Secret: secret,
	}
}
