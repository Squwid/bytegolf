package main

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Squwid/bytegolf/aws"
	"github.com/Squwid/bytegolf/runner"
)

// motd returns a string which is the message of the day
// todo: come back here and add more messages
func motd() string {
	var messages = []string{"Bugs not birds..."}
	// random := rand.Intn(len(messages) - 1)
	return messages[0]
}

func index(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", struct {
		MOTD string
	}{
		MOTD: motd(),
	})
}

func account(w http.ResponseWriter, req *http.Request) {
	if !loggedIn(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "account.html", struct{}{})
}

func leaderboards(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "leaderboards.html", struct{}{})
}

func login(w http.ResponseWriter, req *http.Request) {
	if loggedIn(w, req) {
		http.Redirect(w, req, "/account", http.StatusSeeOther)
		return
	}

	// if the user is trying to login
	if req.Method == "POST" {
		reqEmail := req.FormValue("login_email")
		reqPass := req.FormValue("login_password")
		correctLogin, err := tryLogin(reqEmail, reqPass)
		if err != nil {
			logger.Printf("error logging in: %v\n", err)
			http.Error(w, "an internal server error occurred", http.StatusInternalServerError)
			return
		}
		if !correctLogin {
			tpl.ExecuteTemplate(w, "login.html", struct {
				IncorrectPassword bool
			}{
				IncorrectPassword: true,
			})
			return
		}
		// the password is correct
		// their cookie does not exist correctly at this point
		logger.Println("logged in successfully")
		_, err = logOn(w, reqEmail)
		if err != nil {
			logger.Fatalf("error loggin user on %v\n", err)
			return
		}
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "login.html", struct {
		IncorrectPassword bool
	}{
		IncorrectPassword: false,
	})

}

func holes(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "holes.html", struct {
		Questions map[int]aws.Question
	}{
		Questions: questions,
	})
}

func play(w http.ResponseWriter, req *http.Request) {
	intErr := func() {
		http.Error(w, "an internal server error occurred", http.StatusInternalServerError)
	}
	hole := strings.TrimLeft(req.URL.Path, "/play/")
	question, err := getHoleByLink(hole)
	if err != nil {
		http.Redirect(w, req, "/holes", http.StatusSeeOther)
		return
	}
	if req.Method == http.MethodPost {
		if !loggedIn(w, req) {
			http.Redirect(w, req, "/login", http.StatusSeeOther)
			return
		}
		user, err := FetchUser(w, req)
		if err != nil {
			logger.Printf("error fetching user: %v\n", err)
			intErr()
			return
		}

		lang := req.FormValue("language")
		file, fileHead, err := req.FormFile("codefile")
		if err != nil {
			logger.Printf("an error occurred opening a form file: %v\n", err)
			intErr()
			return
		}
		defer file.Close()

		logger.Printf("submission %s %s %s\n", user, hole, lang)

		bs, err := ioutil.ReadAll(file) // todo: swap this to a buffer maybe
		if err != nil {
			logger.Printf("error reading all file : %v\n", err.Error())
			intErr()
			return
		}
		// 		// run the code from the input through the submission system
		runnerClient := runner.NewClient()
		runnerConfig := runner.NewConfiguration(true, true)
		submission := runner.NewCodeSubmission(user.Email, hole, fileHead.Filename, lang, string(bs), runnerClient, runnerConfig)
		runnerResp, err := submission.Send()
		if err != nil {
			logger.Println(err.Error())
			intErr()
			return
		}
		// TODO: Check output file and scoring system
		// the users submission is wrong
		//todo: handle this
		_ = submission
		_ = runnerResp
	}

	// the question exists and was grabbed
	tpl.ExecuteTemplate(w, "play.html", struct {
		Question aws.Question
	}{
		Question: *question,
	})
}
