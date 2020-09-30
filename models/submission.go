package models

// Submission is the submission coming in from the logged in user
type Submission struct {
	HoleID   string `json:"holeId"`
	Code     string `json:"code"`
	Language string `json:"language"`
}

// CompileInput is the input of Jdoodle minus the ClientID and ClientSecret
type CompileInput struct {
	Code         string `json:"script"`
	StdIn        string `json:"StdIn"`
	Language     string `json:"language"`
	VersionIndex string `json:"versionIndex"`
}

type CompileOutput struct {
	Output     string `json:"output"`
	StatusCode int    `json:"statusCode"`
	Memory     string `json:"memory"`
	CPUTime    string `json:"cpuTime"`
}
