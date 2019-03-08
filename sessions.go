package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/Squwid/bytegolf/aws"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func logOn(w http.ResponseWriter, email string) (aws.User, error) {
	// getting the user first to make sure that it doesnt error out after putting the user in the map
	user, err := getAwsUser(email)
	if err != nil {
		return user, err
	}
	id, err := uuid.NewV4()
	if err != nil {
		return user, err // this code is probably not accessable but just in case of bs
	}

	idString := id.String()
	cookie := &http.Cookie{
		Name:  "bgsession",
		Value: idString,
	}
	http.SetCookie(w, cookie)
	sessions[idString] = session{
		Email:        email,
		lastActivity: time.Now(),
	}
	return user, nil
}

// tryLogin tries an email and password and checks to see if its correct. It uses user caching
// incase the user tries multiple logins. Returns errors if aws does not act as intended
func tryLogin(email, password string) (bool, error) {
	user, err := getAwsUser(email)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logger.Printf("%s tried to login incorrectly\n", email)
		return false, nil
	}
	return true, nil
}

// FetchUser checks to see if the user is logged in, and redirects them to login if they are not logged in. It also checks
// to make sure that the pointer to the user is not nil, and if it is an error is sent back
func FetchUser(w http.ResponseWriter, req *http.Request) (aws.User, error) {
	if !loggedIn(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return aws.User{}, errors.New("this user is not logged in and has been redirected")
	}
	cookie, err := req.Cookie("bgsession")
	if err != nil {
		return aws.User{}, err
	}
	if sess, ok := sessions[cookie.Value]; ok {
		user, err := getAwsUser(sess.Email)
		return user, err
	}
	return aws.User{}, errors.New("should be unreachable code")
}

// loggedIn checks to see if a player is currently logged in
func loggedIn(w http.ResponseWriter, req *http.Request) bool {
	cookie, err := req.Cookie("bgsession")
	if err != nil {
		return false
	}
	// fmt.Printf("SESSIONS: %v\tMY COOKIEVAL: %s\n", sessions, cookie.Value)
	var ok bool
	if session, ok := sessions[cookie.Value]; ok {
		session.lastActivity = time.Now()
		return true
	}
	return ok
}
