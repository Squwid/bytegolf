package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Squwid/bytegolf/globals"
	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

// LoginHandler will send the request to Github to make sure that the user is logged in
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	l := log.WithField("Action", "LoginRequest")
	l.Infof("New login request")

	// Check if user is already logged in
	claims := LoggedIn(r)
	if claims != nil {
		l.Infof("Already logged in")
		http.Redirect(w, r, globals.FrontendAddr()+"/profile/"+claims.BGID, http.StatusSeeOther)
		return
	}

	// Create the github request for the upcoming redirect
	ghReq, err := http.NewRequest("GET", "https://github.com/login/oauth/authorize", nil)
	if err != nil {
		l.WithError(err).Errorf("Error creating new request for Github")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Query Strings
	q := ghReq.URL.Query()
	q.Add("client_id", githubClient)
	q.Add("state", githubState)
	q.Add("allow_signup", "true")

	ghReq.URL.RawQuery = q.Encode()
	redirectTo := "https://github.com" + ghReq.URL.RequestURI()
	l.WithField("Redirect", redirectTo).Debugf("Redirect to github")

	// Redirect using the Github URL
	http.Redirect(w, ghReq, redirectTo, http.StatusSeeOther)
}

// CallbackHandler is the callback from Github to grab the auth token
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	l := log.WithField("Action", "CallbackResponse")

	// Get code and state from response body
	codeResp := r.URL.Query().Get("code")
	stateResp := r.URL.Query().Get("state")

	l.WithFields(log.Fields{
		"Code":  codeResp,
		"State": stateResp,
	}).Infof("Github callback")

	// Check state
	if stateResp != githubState {
		l.Warnf("State %s does not match expected state", stateResp)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Call github back
	body, _ := json.Marshal(struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code"`
		State        string `json:"state"`
	}{
		ClientID:     githubClient,
		ClientSecret: githubSecret,
		Code:         codeResp,
		State:        stateResp,
	})

	// Send post request
	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewReader(body))
	if err != nil {
		l.WithError(err).Errorf("Error creating post request to swap code for access token")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Send Request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		l.WithError(err).Errorf("Error sending request to Github")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Make sure status code was 200
	if resp.StatusCode != 200 {
		l.Errorf("Invalid status code back from Github: %v (%v)", resp.Status, resp.StatusCode)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Parse access token from Github
	var authResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		l.WithError(err).Errorf("Error decoding access_token resp from Github: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Make sure access token exists
	if authResp.AccessToken == "" {
		l.Errorf("Expected auth token but it was blank") // Gets called if code is invalid
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	l.Infof("Got access token from Github")

	// Use Access Token to get User
	githubUser, err := GetGithubUser(authResp.AccessToken)
	if err != nil {
		l.WithError(err).Errorf("Error getting Github user")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get Bytegolf User from GithubUser
	bgUser, err := Bytegolf(githubUser)
	if err != nil {
		l.WithError(err).Errorf("Error swapping Github User -> Bytegolf User")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// TODO: Move JWT logic somewhere else
	timeoutDur := time.Hour * 48
	expires := time.Now().Add(timeoutDur)

	// Claims
	claims := Claims{
		BGID: bgUser.BGID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Create a new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token using key
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		l.WithError(err).Errorf("Error signing token string")
		http.Error(w, "Invalid signing token", http.StatusInternalServerError)
		return
	}

	// Set the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    signedToken,
		Expires:  expires,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	// Successful, redirect
	http.Redirect(w, r, globals.FrontendAddr()+"/profile", http.StatusSeeOther)
}

// TODO: Should the method of getting the token change?
func JWT(token *jwt.Token) (interface{}, error) {
	return jwtKey, nil
}

func LoggedIn(r *http.Request) *Claims {
	claims, ok := r.Context().Value("Claims").(*Claims)
	if !ok {
		return nil
	}
	return claims
}
