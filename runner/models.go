package runner

import (
	"fmt"
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
	Stdin        string `json:"stdin"`

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
	Username   string `json:"username"`
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

// LbOverview holds the data for one spot on the leaderboard for each hole (just the overview version)
type LbOverview struct {
	Username string `json:"username"`
	Language string `json:"language"`
	Score    int    `json:"score"`
}

// Client is the client that allows the user to talk to the API
type Client struct {
	ID     string `json:"clientId"`
	Secret string `json:"clientSecret"`
}

// PrevAnswered holds the information about a users previously answered question
type PrevAnswered struct {
	Correct  bool   `json:"correct"`
	Language string `json:"lanugage"`
	Score    int    `json:"score"`
}

// NewClient returns a new client using the environmental variables
// of RUNNER_ID for the ID and RUNNER_SECRET for the secret
func NewClient() *Client {
	fmt.Println(os.Getenv("RUNNER_ID"), os.Getenv("RUNNER_SECRET"))
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
func NewCodeSubmission(email, username, questionID, input, filename, language, code string, client *Client, sess *session.Session) *CodeSubmission {
	id, _ := uuid.NewV4()
	return &CodeSubmission{
		UUID:         id.String(),
		Script:       code,
		Language:     language,
		VersionIndex: "0",
		ID:           client.ID,
		Secret:       client.Secret,
		Stdin:        input,
		Info: &FileInfo{
			QuestionID: questionID,
			Name:       filename,
			User:       email,
			Username:   username,
		},
		awsSess: sess,
	}
}
