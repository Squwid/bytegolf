package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Squwid/bytegolf/aws"
	uuid "github.com/satori/go.uuid"
)

func setUserHole(hole int, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:  "hole",
		Value: strconv.Itoa(hole),
	})
}

func getUserHole(req *http.Request) int {
	var cookie *http.Cookie
	cookie, err := req.Cookie("hole")
	if err != nil {
		// user either deleted cookie or they havent joined a game yet
		return 1
	}
	i, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return 1
	}
	if i > 9 {
		return 9
	}
	if i < 1 {
		return 1
	}
	return i
}

func getUser(req *http.Request) (*aws.User, error) {
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
			return u, nil
		}
	}
	u, err := getAwsUser(sessions[cookie.Value].Email)
	// u, err := aws.GetUser(sessions[cookie.Value].Email)
	if err != nil {
		return nil, err
	}
	users = append(users, u)
	return u, nil
}

func loggedIn(req *http.Request) bool {
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
