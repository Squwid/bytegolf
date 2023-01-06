package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ErrBadGithubStatus gets returned from GetGithubUser if a 200
// is not returned from github when making calls.
var ErrBadGithubStatus = errors.New("bad status code from github")

// fetchUserFromGithub gets a github user using their access token.
func fetchUserFromGithub(token string) (*GithubUser, error) {
	r, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	r.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrBadGithubStatus
	}

	var ghu GithubUser
	if err := json.NewDecoder(resp.Body).Decode(&ghu); err != nil {
		return nil, err
	}
	return &ghu, nil
}

// callback returns the access token from Github if the code
// and state are correct.
func callback(code, state string) (*string, error) {
	bs, _ := json.Marshal(struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code"`
		State        string `json:"state"`
	}{
		ClientID:     githubClient,
		ClientSecret: githubSecret,
		Code:         code,
		State:        state,
	})
	req, err := http.NewRequest("POST",
		"https://github.com/login/oauth/access_token", bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, ErrBadGithubStatus
	}

	var authResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	return &authResp.AccessToken, nil
}
