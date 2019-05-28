package main

// variables that are needed for the database
var (
	dbUsername string
	dbPassword string
	dbHost     string
	dbName     string
)

// initializeDB initializes the database using connection strings
func initializeDB() {

}

// GithubUser is the user structure that uses github
type GithubUser struct {
	ID         int    `json:"id"`
	Username   string `json:"login"`
	PictureURI string `json:"avatar_url"`
	GithubURI  string `json:"html_url"`
	Name       string `json:"name"`
}

// BGUser is the type that gets stored in the database and used for the profile page
type BGUser struct {
	GithubUser       *GithubUser `json:"githubUser"`
	TotalSubmissions int         `json:"totalSubmissions"`
	TotalCorrect     int         `json:"totalCorrect"`
	TotalBytes       int         `json:"totalBytes"`
}

// NewBGUser creates a new BytegolfUser from the Github User
func NewBGUser(ghu *GithubUser) *BGUser {
	// TODO: this needs to grab the user from the database rather than just create a new one
	return &BGUser{
		GithubUser:       ghu,
		TotalSubmissions: 0,
		TotalCorrect:     0,
		TotalBytes:       0,
	}
}

// Exists checks to see if a github user already exists in the database
func (ghu *GithubUser) Exists() (bool, error) {
	return false, nil
}
