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
		if u.Email == sessions[cookie.Value].Email {
			return u
		}
	}
	u, err := aws.GetUser(sessions[cookie.Value].Email)
	if err != nil {
		panic(err)
	}
	users = append(users, u)
	return u
}

func loggedIn(w http.ResponseWriter, req *http.Request) bool {
	cookie, err := req.Cookie("session")
	if err != nil {
		return false
	}
	var ok bool
	if session, ok := sessions[cookie.Value]; ok {
		session.lastActivity = time.Now()
		sessions[cookie.Value] = session
	}
	return ok
}
