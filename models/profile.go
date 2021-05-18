package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

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
