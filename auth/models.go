package auth

import (
	"os"
	"time"

	"github.com/Squwid/bytegolf/models"
	"github.com/google/uuid"
)

const cookieName = "bg-token"
const loginRedirect = "/" // Where to go post login

var githubClient, githubSecret, githubState string
var jwtKey []byte

func init() {
	githubClient = os.Getenv("GITHUB_CLIENT")
	githubSecret = os.Getenv("GITHUB_SECRET")
	githubState = os.Getenv("GITHUB_STATE")
	jwtKey = []byte(os.Getenv("JWT_SECRET"))

	if githubClient == "" || githubSecret == "" || githubState == "" || len(jwtKey) == 0 {
		panic("missing github env variables")
	}
}

// NewBytegolfUser returns a new bytegolf user
func NewBytegolfUser(ghu models.GithubUser) *models.BytegolfUser {
	return &models.BytegolfUser{
		BGID:        uuid.New().String(),
		GithubUser:  ghu,
		CreatedTime: time.Now().UTC(),
	}
}
