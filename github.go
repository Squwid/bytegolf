package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Session information regarding github logins
var (
	sessions    = map[string]*session{}
	sessionLock = &sync.RWMutex{}
)

type session struct {
	User         *GithubUser
	lastActivity time.Time
}

// github variables
var (
	clientID     string
	clientSecret string
)

func setGitClient() {
	type oauth struct {
		ClientID     string `json:"ClientID"`
		ClientSecret string `json:"ClientSecret"`
	}

	dat, err := ioutil.ReadFile("github_oauth.json")
	if err != nil {
		panic(err)
	}

	var oa oauth
	err = json.Unmarshal(dat, &oa)
	if err != nil {
		panic(err)
	}
	clientID = oa.ClientID
	clientSecret = oa.ClientSecret
}

func githubOAUTH(w http.ResponseWriter, req *http.Request) {
	// get the github response back to the site and make sure it matches everything from before
	code := req.URL.Query().Get("code")
	state := req.URL.Query().Get("state")

	type params struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code"`
		State        string `json:"state"`
	}

	var p = &params{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Code:         code,
		State:        state,
	}

	// TODO: this came from postman but needs to change
	s := fmt.Sprintf("------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"client_id\"\r\n\r\n%s\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"client_secret\"\r\n\r\n%s\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"code\"\r\n\r\n%s\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"state\"\r\n\r\n%s\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--", p.ClientID, p.ClientSecret, p.Code, p.State)

	payload := strings.NewReader(s)

	newReq, err := http.NewRequest(http.MethodPost, "https://github.com/login/oauth/access_token", payload)
	if err != nil {
		logger.Printf("error creating request: %v\n", err)
		return
	}

	newReq.Header.Add("accept", "application/json")
	newReq.Header.Add("content-type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	newReq.Header.Add("cache-control", "no-cache")

	resp, err := http.DefaultClient.Do(newReq)
	if err != nil {
		logger.Printf("error sending request: %v\n", err)
		return
	}

	type ghresp struct {
		AccessToken string `json:"access_token"`
	}

	defer resp.Body.Close()
	respBS, _ := ioutil.ReadAll(resp.Body)

	var ghr ghresp
	err = json.Unmarshal(respBS, &ghr)
	if err != nil {
		logger.Printf("error unmarshalling response: %v\n", err)
		return
	}

	// get the user using the login token
	userReq, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		logger.Fatalf("error creating user request: %v\n", err)
		return
	}
	userReq.Header.Add("authorization", fmt.Sprintf("Bearer %s", ghr.AccessToken))
	userResp, err := http.DefaultClient.Do(userReq)
	if err != nil {
		logger.Fatalf("error sending user request: %v\n", err)
		return
	}

	defer userResp.Body.Close()
	userBS, err := ioutil.ReadAll(userResp.Body)
	if err != nil {
		logger.Fatalf("error reading userBS: %v\n", err)
		return
	}

	// get the github user and add them to the session
	var githubUser GithubUser
	err = json.Unmarshal(userBS, &githubUser)
	if err != nil {
		logger.Fatalf("error marshalling github user: %v\n", err)
		return
	}
	githubUser.login(w, req)

	// redirect the user to the main page after they sign in
	// tpl.ExecuteTemplate(w, "profile.html", githubUser)
	http.Redirect(w, req, "/account", http.StatusSeeOther)
}

func gitLogin(w http.ResponseWriter, req *http.Request) {
	// the user is already logged in so dont send them to the auth page
	if loggedIn(w, req) {
		http.Redirect(w, req, "/account", http.StatusSeeOther)
		return
	}
	logger.Printf("logging a user into github...\n")

	newReq, err := http.NewRequest(http.MethodGet, "https://github.com/login/oauth/authorize", nil)
	if err != nil {
		logger.Fatalln(err)
		return
	}

	q := newReq.URL.Query()
	q.Add("client_id", clientID)
	q.Add("state", gitState)
	q.Add("allow_signup", "true")
	newReq.URL.RawQuery = q.Encode()
	http.Redirect(w, newReq, "https://github.com"+newReq.URL.RequestURI(), http.StatusSeeOther)
}
