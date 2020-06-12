package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// ErrBadGithubStatus gets returned from GetGithubUser if a 200 is not returned from github
var ErrBadGithubStatus = errors.New("bad status code from github")

// GetGithubUser gets a github user using their access token
func GetGithubUser(token string) (*GithubUser, error) {
	// Create request
	r, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	// Attach token header
	r.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	// Send request
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Warnf("Expected 200 from github but got %s (%v)", resp.Status, resp.StatusCode)
		return nil, ErrBadGithubStatus
	}

	// Parse to github object
	var ghu GithubUser
	if err := json.NewDecoder(resp.Body).Decode(&ghu); err != nil {
		return nil, err
	}

	return &ghu, nil
}

// ProfileHandler is the handler to display a user's profile
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn, claims := LoggedIn(r)
	if !loggedIn {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"LoggedIn": false}`))
		return
	}

	// Get user from BGID
	user, err := bytegolfUserFromBGID(claims.BGID)
	if err != nil {
		log.WithField("BGID", claims.BGID).WithError(err).Errorf("Got error getting users profile")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		log.WithField("BGID", claims.BGID).Warnf("Expected user, but was not found")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Got User, Parse & Show
	bs, err := json.Marshal(user)
	if err != nil {
		log.WithError(err).Errorf("Error marshalling claims")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(bs)
}
