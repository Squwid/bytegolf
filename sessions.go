package main

import (
	"net/http"
	"time"

	"github.com/Squwid/bytegolf/bgaws"
	uuid "github.com/satori/go.uuid"
)

func getUser(w http.ResponseWriter, req *http.Request) *bgaws.User {
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
		if u.Username == currentSessions[cookie.Value].Username {
			return u
		}
	}
	u, _ := bgaws.GetUser(currentSessions[cookie.Value].Username) //TODO: when does this error get called?
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
		currentSessions[cookie.Value] = session // todo: does this do what i think it does?
	}
	// refresh session
	// c.MaxAge = sessionLength
	// http.SetCookie(w, cookie)
	return ok
}
