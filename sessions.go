package main

import (
	"net/http"
	"strconv"
	"time"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func fetchUser(w http.ResponseWriter, req *http.Request) (*GithubUser, error) {
	if !loggedIn(w, req) {
		return nil, nil
	}

	cookie, err := req.Cookie("bgsession")
	if err != nil {
		return nil, err
	}

	sessionLock.RLock()
	defer sessionLock.RUnlock()
	return sessions[cookie.Value].User, nil
}

func loggedIn(w http.ResponseWriter, req *http.Request) bool {
	cookie, err := req.Cookie("bgsession")
	if err != nil {
		return false
	}

	var ok bool
	sessionLock.RLock()
	defer sessionLock.RUnlock()
	if session, ok := sessions[cookie.Value]; ok {
		session.lastActivity = time.Now()
		return true
	}
	return ok
}

// this will login a user, make sure they are not already logged in first
func (user *GithubUser) login(w http.ResponseWriter, req *http.Request) {
	// hash the github login id to not store it raw
	uid, err := bcrypt.GenerateFromPassword([]byte(strconv.Itoa(user.ID)), bcrypt.MinCost)
	if err != nil {
		logger.Printf("error generating cookie: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:  "bgsession",
		Value: string(uid),
		Path: "/",
	}

	http.SetCookie(w, cookie)
	req.AddCookie(cookie)

	s := &session{
		User:         user,
		lastActivity: time.Now(),
	}

	// add the session to the list of sessions
	sessionLock.Lock()
	sessions[string(uid)] = s
	sessionLock.Unlock()

	// run a function to remove the session in a day
	go func(id string) {
		time.Sleep(time.Hour * 24)
		sessionLock.Lock()
		delete(sessions, id)
		sessionLock.Unlock()
	}(string(uid))
}
