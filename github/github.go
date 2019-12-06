package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Squwid/bytegolf/sess"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var state = "abcdefg"

// Login handler
func Login(w http.ResponseWriter, req *http.Request) {
	loggedIn, _, err := sess.LoggedIn(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("error logging in: %v", err)
		return
	}
	if loggedIn {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		log.Info("user tried to log in but they already are logged in")
		return
	}

	log.Info("logging a new user into github")
	req2, err := http.NewRequest(http.MethodGet, "https://github.com/login/oauth/authorize", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("error creating new github get request: %v", err)
		return
	}

	gitClient := getClient()
	q := req2.URL.Query()
	q.Add("client_id", gitClient.ID)
	q.Add("state", state)
	q.Add("allow_signup", "true")
	req2.URL.RawQuery = q.Encode()
	http.Redirect(w, req2, "https://github.com"+req2.URL.RequestURI(), http.StatusSeeOther)
}

// Oauth handler
func Oauth(w http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")
	state := req.URL.Query().Get("state")
	c := getClient()
	type params struct {
		client
		Code  string `json:"code"`
		State string `json:"state"`
	}

	var p = &params{
		client: c,
		Code:   code,
		State:  state,
	}
	s := fmt.Sprintf("------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"client_id\"\r\n\r\n%s\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"client_secret\"\r\n\r\n%s\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"code\"\r\n\r\n%s\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"state\"\r\n\r\n%s\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--", p.client.ID, p.client.Secret, p.Code, p.State)

	payload := strings.NewReader(s)

	newReq, err := http.NewRequest(http.MethodPost, "https://github.com/login/oauth/access_token", payload)
	if err != nil {
		log.Printf("error creating request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newReq.Header.Add("accept", "application/json")
	newReq.Header.Add("content-type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	newReq.Header.Add("cache-control", "no-cache")

	resp, err := http.DefaultClient.Do(newReq)
	if err != nil {
		log.Printf("error sending request: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
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
		log.Printf("error unmarshalling response: %v\n", err)
		return
	}

	// get the user using the login token
	userReq, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		log.Fatalf("error creating user request: %v\n", err)
		return
	}
	userReq.Header.Add("authorization", fmt.Sprintf("Bearer %s", ghr.AccessToken))
	userResp, err := http.DefaultClient.Do(userReq)
	if err != nil {
		log.Fatalf("error sending user request: %v\n", err)
		return
	}

	defer userResp.Body.Close()
	userBS, err := ioutil.ReadAll(userResp.Body)
	if err != nil {
		log.Fatalf("error reading userBS: %v\n", err)
		return
	}

	// get the github user and add them to the session
	var githubUser User
	err = json.Unmarshal(userBS, &githubUser)
	if err != nil {
		log.Fatalf("error marshalling github user: %v\n", err)
		return
	}
	githubUser.login(w, req)

	// redirect the user to the main page after they sign in
	// tpl.ExecuteTemplate(w, "profile.html", githubUser)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (user *User) login(w http.ResponseWriter, req *http.Request) {
	// hash the github login id to not store it raw

	// grab the bgid from the github user in storage and attach it to a session
	// for the user to be logged in
	var update bool
	u, err := getUserFromGithub(user.GithubID)
	if err == ErrNotFound {
		// the user was not found lets create him
		uid := uuid.New().String()
		u = &User{
			BGID: uid,
		}
		update = true // update the new user to user tables
	} else if err != nil {
		log.Errorf("error logging in github user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if u == nil {
		log.Errorf("user is nil for some reason: %v", user.GithubID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: add logic here to check for github name changes and update the table accordingly
	// copy over values to the user object
	user.BGID = u.BGID
	if user.Username != u.Username {
		update = true
	}

	// TODO: this uuid is getting created on every login but we need to store github users once,
	// and then grab the old one

	s, err := sess.Login(user.BGID)
	if err != nil {
		log.Fatalf("error logging in after github: %v", err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "bg-session",
		Value: s.ID,
		Path:  "/",
	})
	if update {
		err = user.Put()
		if err != nil {
			log.Errorf("error updating user %s: %v", user.BGID, err)
		}
	}
}
