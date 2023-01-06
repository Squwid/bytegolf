package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/uptrace/bun"
)

// CookieName is the key of the cookie for the environment.
var CookieName string

var githubClient, githubSecret, githubState string
var jwtKey []byte

func init() {
	githubClient = os.Getenv("GITHUB_CLIENT")
	githubSecret = os.Getenv("GITHUB_SECRET")
	githubState = os.Getenv("GITHUB_STATE")
	jwtKey = []byte(os.Getenv("JWT_SECRET"))
	CookieName = os.Getenv("BG_COOKIE_NAME")

	if githubClient == "" || githubSecret == "" ||
		githubState == "" || len(jwtKey) == 0 || CookieName == "" {
		panic("missing github env variables")
	}
}

// NewBytegolfUser returns a new bytegolf user
func NewBytegolfUser(ghu GithubUser) *BytegolfUserDB {
	return &BytegolfUserDB{
		// TODO: Randomize a custom BGID.
		GithubUser:      ghu,
		BGID:            fmt.Sprintf("%v", ghu.GithubID),
		LastUpdatedTime: time.Now().UTC(),
		CreatedTime:     time.Now().UTC(),
	}
}

// GithubUser gets returned from Github and is composed by BytegolfUser.
type GithubUser struct {
	GithubID  int64     `json:"id" bun:"id,pk,notnull"`
	Login     string    `json:"login" bun:"login"`
	URL       string    `json:"html_url" bun:"github_url"`
	AvatarURL string    `json:"avatar_url" bun:"github_avatar_url"`
	UpdatedAt time.Time `json:"updated_at" bun:"-"`
}

// BytegolfUserDB is the database object for the 'users' table.
type BytegolfUserDB struct {
	bun.BaseModel `bun:"table:users"`

	BGID            string    `bun:"bgid,notnull"`
	LastUpdatedTime time.Time `bun:"updated_time,notnull"`
	CreatedTime     time.Time `bun:"created_time,notnull"`

	GithubUser
}

// BytegolfUserClient is the object that gets returned to the client when making
// profile calls.
type BytegolfUserClient struct {
	GithubID    string `json:"GithubId"`
	BGID        string `json:"BGID"`
	DisplayName string `json:"DisplayName"`
	GithubURL   string `json:"GithubUrl"`
	AvatarURL   string `json:"AvatarUrl"`
}

func (bgdb BytegolfUserDB) ToProfile() BytegolfUserClient {
	return BytegolfUserClient{
		GithubID:    fmt.Sprintf("%v", bgdb.GithubID),
		BGID:        bgdb.BGID,
		DisplayName: bgdb.GithubUser.Login,
		GithubURL:   bgdb.GithubUser.URL,
		AvatarURL:   bgdb.GithubUser.AvatarURL,
	}
}

// Claims = JWTClaims.
type Claims struct {
	GithubID int64  `json:"GithubId"`
	BGID     string `json:"BGID"`

	jwt.StandardClaims
}
