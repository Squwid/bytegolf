package main

import (
	"fmt"
	"net/http"

	"github.com/Squwid/bytegolf/github"
	"github.com/Squwid/bytegolf/sess"
	log "github.com/sirupsen/logrus"
)

// isLoggedIn is the api call that the front end makes to see if a user is signed in
func isLoggedIn(w http.ResponseWriter, req *http.Request) {
	l := log.WithField("Action", "IsLoggedIn")

	w.Header().Set("Content-Type", "application/json")
	loggedIn, s, err := sess.LoggedIn(req)
	if err != nil {
		l.Errorf("Error checking if logged in: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !loggedIn {
		l.Infof("User is not logged in")
		w.Write([]byte(`{"logged_in": false}`))
		return
	}

	l = l.WithField("BGID", s.BGID)

	// Get user's info to get their user info
	user, err := github.RetreiveUser(s.BGID)
	if err != nil {
		l.Errorf("Error retreiving github user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	l.WithField("Username", user.Username).Infof("User is logged in")
	w.Write([]byte(fmt.Sprintf(`{"logged_in": true, "username": "%s"}`, user.Username)))
}
