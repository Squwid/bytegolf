package auth

import (
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Squwid/bytegolf/db"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const CookieName = "bg-token"

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
func NewBytegolfUser(ghu GithubUser) *BytegolfUser {
	return &BytegolfUser{
		BGID:        uuid.New().String(),
		GithubUser:  ghu,
		CreatedTime: time.Now().UTC(),
	}
}

// GithubUser is the object that comes back from github on a user lookup
type GithubUser struct {
	ID        int64     `json:"id"`
	Login     string    `json:"login"`
	URL       string    `json:"html_url"`
	AvatarURL string    `json:"avatar_url"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BytegolfUser is the structure of how bytegolf user's are stored in the database
type BytegolfUser struct {
	BGID            string
	LastUpdatedTime time.Time
	CreatedTime     time.Time

	GithubUser GithubUser
}

// BytegolfUserProfile is the BytegolfUser struct but with no sensitive fields
type BytegolfUserProfile struct {
	BGID        string `json:"BGID"`
	DisplayName string `json:"DisplayName"`
	GithubURL   string `json:"GithubUrl"`
	AvatarURL   string `json:"AvatarUrl"`
}

// ToProfile takes the BytegolfUser (Database object) and mutates it to a Profile (Frontend object)
func (bgu BytegolfUser) ToProfile() BytegolfUserProfile {
	return BytegolfUserProfile{
		BGID:        bgu.BGID,
		DisplayName: bgu.GithubUser.Login,
		GithubURL:   bgu.GithubUser.URL,
		AvatarURL:   bgu.GithubUser.AvatarURL,
	}
}

// Claims is what gets stored in the JWT
type Claims struct {
	BGID string `json:"BGID"`

	jwt.StandardClaims
}

func (bgu *BytegolfUser) Collection() *firestore.CollectionRef { return db.ProfileCollection() }
func (bgu *BytegolfUser) DocID() string                        { return bgu.BGID }
func (bgu *BytegolfUser) Data() interface{} {
	bgu.LastUpdatedTime = time.Now().UTC()
	return *bgu
}
