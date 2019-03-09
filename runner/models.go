package runner

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	uuid "github.com/satori/go.uuid"
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

// CodeFile holds both the response and the submission along with information about the question
type CodeFile struct {
	Submission CodeSubmission `json:"submission"`
	Response   CodeResponse   `json:"response"`

	Correct bool `json:"correct"`
	Length  int  `json:"length"`
}

// CodeSubmission is what gets submitted to the
type CodeSubmission struct {
	UUID         string `json:"uuid"`
	Script       string `json:"script"`
	Language     string `json:"language"`
	VersionIndex string `json:"versionIndex"`

	ID     string    `json:"clientId"`
	Secret string    `json:"clientSecret"`
	Info   *FileInfo `json:"info"`

	// handles aws s3 storage
	awsSess *session.Session
}

// FileInfo todo
type FileInfo struct {
	QuestionID string `json:"questionID"`
	Name       string `json:"name"`
	User       string `json:"user"`
}

// CodeResponse is the response from the Code Runner API that gets a result
type CodeResponse struct {
	UUID       string    `json:"uuid"`
	Output     string    `json:"output"`
	StatusCode int       `json:"statusCode"`
	Memory     string    `json:"memory"`
	CPUTime    string    `json:"cpuTime"`
	Info       *FileInfo `json:"info"`

	// Information regarding the response post check
	Correct bool `json:"correct"`
	Length  int  `json:"length"`

	awsSess *session.Session
}

// Client is the client that allows the user to talk to the API
type Client struct {
	ID     string `json:"clientId"`
	Secret string `json:"clientSecret"`
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

// NewCodeSubmission todo:
func NewCodeSubmission(username, questionID, filename, language, code string, client *Client, sess *session.Session) *CodeSubmission {
	id, _ := uuid.NewV4()
	return &CodeSubmission{
		UUID:         id.String(),
		Script:       code,
		Language:     language,
		VersionIndex: "0",
		ID:           client.ID,
		Secret:       client.Secret,
		Info: &FileInfo{
			QuestionID: questionID,
			Name:       filename,
			User:       username,
		},
		awsSess: sess,
	}
}
