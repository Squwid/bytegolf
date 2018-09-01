package runner

import (
	"os"

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

// CodeSubmission is what gets submitted to the
type CodeSubmission struct {
	UUID         string `json:"uuid"`
	Script       string `json:"script"`
	Language     string `json:"language"`
	VersionIndex string `json:"versionIndex"`
	ID           string `json:"clientId"`
	Secret       string `json:"clientSecret"`

	Info   *FileInfo      `json:"info"`
	Config *Configuration `json:"config"`
}

// Configuration holds information on where to save the submissions and responses to
// and whether or not it should store logs
type Configuration struct {
	SaveLogs        bool `json:"saveLogs"`
	SaveSubmissions bool `json:"saveSubmissions"`
}

// FileInfo todo
type FileInfo struct {
	FileName string `json:"fileName"`
	User     string `json:"user"`
	Game     string `json:"game"` // TODO: This was changed in the package to key instead of game
	GameName string `json:"gameName"`
}

// CodeResponse is the response from the Code Runner API that gets a result
type CodeResponse struct {
	UUID       string `json:"uuid"`
	Output     string `json:"output"`
	StatusCode int    `json:"statusCode"`
	Memory     string `json:"memory"`
	CPUTime    string `json:"cpuTime"`

	Info *FileInfo
}

// Client is the client that allows the user to talk to the API
type Client struct {
	ID     string `json:"clientId"`
	Secret string `json:"clientSecret"`
}

// NewConfiguration returns a pointer to a configuration
func NewConfiguration(saveLogs, saveSubmissions bool) *Configuration {
	return &Configuration{
		SaveLogs:        saveLogs,
		SaveSubmissions: saveSubmissions,
	}

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
func NewCodeSubmission(username, gameName, gameID, filename, language, code string, client *Client, config *Configuration) *CodeSubmission {
	id, _ := uuid.NewV4()
	return &CodeSubmission{
		UUID:         id.String(),
		Script:       code,
		Language:     language,
		VersionIndex: "0",
		ID:           client.ID,
		Secret:       client.Secret,
		Config:       config,
		Info: &FileInfo{
			FileName: filename,
			User:     username,
			Game:     gameID,
			GameName: gameName,
		},
	}
}
