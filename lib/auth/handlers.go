package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Squwid/bytegolf/lib/globals"
	"github.com/Squwid/bytegolf/lib/log"
	"github.com/sirupsen/logrus"
)

// LoginHandler will send the request to Github to make sure that the user is logged in
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLogger().WithField("Action", "LoginRequest")
	logger.Debugf("New login request.")

	claims := LoggedIn(r)
	if claims != nil {
		http.Redirect(w, r, globals.FrontendAddr()+"/profile/"+claims.BGID, http.StatusSeeOther)
		logger.Debugf("User already logged in, redirecting.")
		return
	}

	ghReq, _ := http.NewRequest("GET", "https://github.com/login/oauth/authorize", nil)
	qs := ghReq.URL.Query()
	qs.Add("client_id", githubClient)
	qs.Add("state", githubState)
	qs.Add("allow_signup", "true")
	ghReq.URL.RawQuery = qs.Encode()

	redirectTo := ghReq.URL.String()
	http.Redirect(w, ghReq, redirectTo, http.StatusSeeOther)

	logger.WithField("Redirect", redirectTo).Debugf("Github login redirect.")
}

// CallbackHandler is the callback from Github to grab the auth token
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLogger().WithField("Action", "CallbackResponse")

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	logger.WithFields(logrus.Fields{
		"Code":  code,
		"State": state,
	}).Debugf("Github callback.")

	if state != githubState {
		http.Error(w, "Bad request", http.StatusBadRequest)
		logger.Warnf("State %s does not match expected state.", state)
		return
	}

	accessToken, err := callback(code, state)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		logger.WithError(err).Errorf("Could not call Github back.")
		return
	}

	// Use Access Token to call Github to fetch user.
	githubUser, err := fetchUserFromGithub(*accessToken)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.WithError(err).
			Errorf("Error getting Github user using the access token.")
		return
	}

	// Check if user already exists, create if not.
	ctx := context.Background()
	bgUser, err := createOrGetDBUser(ctx, githubUser)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.WithError(err).
			Errorf("Error swapping Github User -> Bytegolf User.")
		return
	}

	if err := writeJWT(w, bgUser); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.WithError(err).Errorf("Error writing JWT.")
		return
	}
	logger.Debugf("Login flow complete. JWT Written.")

	http.Redirect(w, r, globals.FrontendAddr()+"/profile", http.StatusSeeOther)
}

// ShowClaims shows the claims in the users cookie for the frontend
func ShowClaims(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLogger().WithField("Action", "ShowClaims")

	claims := LoggedIn(r)
	if claims == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"LoggedIn": false}`))
		return
	}

	bs, _ := json.Marshal(claims)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)

	logger.WithFields(logrus.Fields{
		"BGID": claims.BGID,
		"IP":   r.RemoteAddr,
	}).Debugf("Retreived Claims")
}
