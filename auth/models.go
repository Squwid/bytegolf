package auth

import (
	"time"

	"github.com/Squwid/bytegolf/models"
	"github.com/Squwid/bytegolf/secrets"
	"github.com/google/uuid"
)

const profileCollection = "profile"
const cookieName = "bg-token"
const loginRedirect = "/api/profile/"

var state string
var client *secrets.Client
var jwtKey []byte

func init() {
	client = secrets.Must(secrets.GetClient("Github")).(*secrets.Client)

	jwtKey = []byte(secrets.Must(secrets.GetClient("JWT-KEY")).(*secrets.Client).Secret)
	state = secrets.Must(secrets.GetClient("State")).(*secrets.Client).Secret
}

// NewBytegolfUser returns a new bytegolf user
func NewBytegolfUser(ghu models.GithubUser) *models.BytegolfUser {
	return &models.BytegolfUser{
		BGID:        uuid.New().String(),
		Role:        "user",
		GithubUser:  ghu,
		CreatedTime: time.Now().UTC(),
	}
}
