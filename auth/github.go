package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Squwid/bytegolf/models"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// ErrBadGithubStatus gets returned from GetGithubUser if a 200 is not returned from github
var ErrBadGithubStatus = errors.New("bad status code from github")

// GetGithubUser gets a github user using their access token
func GetGithubUser(token string) (*models.GithubUser, error) {
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
		return nil, fmt.Errorf("got bad status code %v", resp.StatusCode)
	}

	// Parse to github object
	var ghu models.GithubUser
	if err := json.NewDecoder(resp.Body).Decode(&ghu); err != nil {
		return nil, err
	}

	return &ghu, nil
}

// ShowClaims shows the claims in the users cookie for the frontend
func ShowClaims(w http.ResponseWriter, r *http.Request) {
	loggedIn, claims := LoggedIn(r)
	if !loggedIn {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"LoggedIn": false}`))
		return
	}

	// Marshal claims and return
	bs, err := json.Marshal(claims)
	if err != nil {
		log.WithError(err).Errorf("Error marshalling claims")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.WithFields(log.Fields{
		"BGID": claims.BGID,
		"IP":   r.RemoteAddr,
	}).Infof("Retreived Claims")
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}

// ShowProfile shows the target user's profile
func ShowProfile(w http.ResponseWriter, r *http.Request) {
	bgid := mux.Vars(r)["bgid"]

	l := log.WithFields(log.Fields{
		"Action": "ShowProfile",
		"BGID":   bgid,
		"IP":     r.RemoteAddr,
	})
	log.Infof("Request to show profile")

	// Get user by BGID
	user, err := bytegolfUserFromBGID(bgid)
	if err != nil {
		l.WithError(err).Errorf("Error getting profile")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// User doesnt exist, return a 404
	if user == nil {
		l.Warnf("Profile not found")
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	// Got user, translate to profile and send to user
	profile := user.ToProfile()

	bs, err := json.Marshal(profile)
	if err != nil {
		l.WithError(err).Errorf("Error marshalling profile")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	l.Infof("Found profile")

	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}
