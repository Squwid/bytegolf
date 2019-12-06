package github

import (
	"context"
	"errors"
	"strings"

	"github.com/Squwid/bytegolf/firestore"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

const collection = "users"

// User is the user structure that uses github
type User struct {
	BGID string `json:"bg_id"`

	GithubID   int    `json:"id"`
	Username   string `json:"login"`
	PictureURI string `json:"avatar_url"`
	GithubURI  string `json:"html_url"`
	Name       string `json:"name"`

	PermissionLevel string `json:"permission_level"`
}

// ErrNotFound is an error that is returned when a doc is not found
var ErrNotFound = errors.New("users not found")

// RetreiveUser exists checks to see if a user exists
func RetreiveUser(id string) (*User, error) {
	ctx := context.Background()
	ref, err := firestore.Client.Collection(collection).Doc(id).Get(ctx)
	if err != nil && strings.Contains(err.Error(), "code = NotFound") {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	var user User
	err = mapstructure.Decode(ref.Data(), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// i need a function that gets a user from the user table by
// using the github id for when the session does not exist (so when they are not)
// logged in
// getUserFromGithub gets a user from the users table with a specific github id
// it will return ErrNotFound if a user is not found
func getUserFromGithub(githubID int) (*User, error) {
	var users = []User{}
	ctx := context.Background()
	iter := firestore.Client.Collection("users").Where("GithubID", "==", githubID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var user User
		err = mapstructure.Decode(doc.Data(), &user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)

	}
	if len(users) == 0 {
		return nil, ErrNotFound
	}
	if len(users) != 1 {
		log.Warnf("user request for %s returned %v results! messy table", githubID, len(users))
	}
	return &users[0], nil
}

// Put creates a new user in firestore
func (user User) Put() error {
	ctx := context.Background()
	_, err := firestore.Client.Collection("users").Doc(user.BGID).Set(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
