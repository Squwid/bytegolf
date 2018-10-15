package main

import (
	"net/http"
	"time"

	"github.com/Squwid/bytegolf/aws"
	uuid "github.com/satori/go.uuid"
)

func getUser(w http.ResponseWriter, req *http.Request) *aws.User {
	var cookie *http.Cookie
	cookie, err := req.Cookie("session")
	if err != nil {
		session, _ := uuid.NewV4()
		cookie = &http.Cookie{
			Name:  "session",
			Value: session.String(),
		}
	}
	// http.SetCookie(w, cookie)

	// if the user exists already, get user
	for _, u := range users {
		if u.Email == currentSessions[cookie.Value].Username {
			return u
		}
	}
	u, err := aws.GetUser(currentSessions[cookie.Value].Username)
	if err != nil {
		panic(err)
	}
	users = append(users, u)
	return u
}

// currentlyLoggedIn checks to see if a user is currently logged in already by
// checking their session ID
func currentlyLoggedIn(w http.ResponseWriter, req *http.Request) bool {
	cookie, err := req.Cookie("session")
	if err != nil {
		return false
	}
	session, ok := currentSessions[cookie.Value]
	if ok {
		// if session currently exists
		session.lastActivity = time.Now()
		currentSessions[cookie.Value] = session
	}
	// refresh session
	// c.MaxAge = sessionLength
	// http.SetCookie(w, cookie)
	return ok
}
