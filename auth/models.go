package auth

import (
	"time"

	"github.com/Squwid/bytegolf/secrets"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const profileCollection = "profile"
const cookieName = "bg-token"
const loginRedirect = "/api/profile/"

var state string
var client *secrets.Client
var jwtKey []byte

func init() {
	client = secrets.Must(secrets.GetClient("BGGH")).(*secrets.Client)

	jwtKey = []byte(secrets.Must(secrets.GetClient("JWT-KEY")).(*secrets.Client).Secret)
	state = secrets.Must(secrets.GetClient("JWT-KEY")).(*secrets.Client).Secret
}

// GithubUser is the object that comes back from github on a user lookup
type GithubUser struct {
	ID        int64     `json:"id"`
	Login     string    `json:"login"`
	URL       string    `json:"html_url"`
	AvatarURL string    `json:"avatar_url"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BytegolfUser is the structure of how bytegolf user's are stored in Github
type BytegolfUser struct {
	BGID            string
	LastUpdatedTime time.Time
	CreatedTime     time.Time

	GithubUser GithubUser
}

// Claims is what gets stored in the JWT
type Claims struct {
	BGID string
	// BytegolfUser BytegolfUser // For debugging

	jwt.StandardClaims
}

// NewBytegolfUser returns a new bytegolf user
func NewBytegolfUser(ghu GithubUser) *BytegolfUser {
	return &BytegolfUser{
		BGID:        uuid.New().String(),
		GithubUser:  ghu,
		CreatedTime: time.Now().UTC(),
	}
}
