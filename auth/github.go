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
		return nil, fmt.Errorf("got bad status code %v", resp.StatusCode)
	}

	// Parse to github object
	var ghu GithubUser
	if err := json.NewDecoder(resp.Body).Decode(&ghu); err != nil {
		return nil, err
	}

	return &ghu, nil
}

// ShowClaims shows the claims in the users cookie for the frontend
func ShowClaims(w http.ResponseWriter, r *http.Request) {
	claims := LoggedIn(r)
	if claims == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
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
	}).Infof("Retreived Claims")

	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}
