package main

import (
	"io/ioutil"
	"net/http"

	"github.com/Squwid/bytegolf/runner"

	"github.com/Squwid/bytegolf/aws"
)

// MessageOfTheDay returns a string that contains message of the day on the index.html page
func MessageOfTheDay() string {
	var messages = []string{"Bugs not birds..."}
	// random := rand.Intn(len(messages) - 1)
	return messages[0]
}

func index(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", struct {
		MOTD string
	}{
		MOTD: MessageOfTheDay(),
	})
}

func account(w http.ResponseWriter, req *http.Request) {
	if !loggedIn(req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "account.html", struct{}{})
}

func leaderboards(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "leaderboards.html", struct{}{})
}

func login(w http.ResponseWriter, req *http.Request) {
	if loggedIn(req) {
		http.Redirect(w, req, "/account", http.StatusSeeOther)
		return
	}

	// if the user is trying to login
	if req.Method == "POST" {
		reqEmail := req.FormValue("email")
		reqPass := req.FormValue("password")
		correctLogin, err := tryLogin(reqEmail, reqPass)
		if err != nil {
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
		//todo: they logged in now what
		_ = correctLogin
	}
	tpl.ExecuteTemplate(w, "login.html", struct {
		IncorrectPassword bool
	}{
		IncorrectPassword: false,
	})

}

func play(w http.ResponseWriter, req *http.Request) {
	internalErr := func(w http.ResponseWriter) {
		http.Error(w, "an internal server error occurred", http.StatusInternalServerError)
	}
	userErr := func(w http.ResponseWriter) {
		http.Error(w, "user misuse error", http.StatusBadRequest)
	}
	var user *aws.User
	if !loggedIn(req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	user, err := getUser(req)
	if err != nil {
		internalErr(w)
		return
	}
	if user == nil {
		internalErr(w)
		logger.Fatalf("this code should be unreachable. NEED TO CHECK ERROR ANYWAYS")
		return
	}

	if !CurrentGame.InProgress() {
		// TODO: here i need to put the create a game option html page
	}

	// if the player has submitted a bytegolf submission
	if req.Method == "POST" {
		lang := req.FormValue("language")
		hole := getUserHole(req)
		file, fileHead, err := req.FormFile("codefile")
		if err != nil {
			logger.Printf("an error occurred opening a form file: %v\n", err)
			userErr(w)
			return
		}
		defer file.Close()
		logger.Printf("%s uploading file to server %s\n", user.Email, fileHead.Filename)
		// read the input file
		bs, err := ioutil.ReadAll(file) // todo: swap to buffer
		if err != nil {
			logger.Println("an error occured reading file", err.Error())
			internalErr(w)
			return
		}

		// run the code from the input through the submission system
		runnerClient := runner.NewClient()
		runnerConfig := runner.NewConfiguration(true, true)
		submission := runner.NewCodeSubmission(user.Email, CurrentGame.Name, CurrentGame.ID, fileHead.Filename, lang, string(bs), runnerClient, runnerConfig)
		runnerResp, err := submission.Send()
		if err != nil {
			logger.Println(err.Error())
			internalErr(w)
			return
		}
		if !CurrentGame.CheckSubmission(hole, runnerResp.Output) {
			// the users submission is wrong
			//todo: handle this
		}
		_ = submission
		_ = runnerResp

	}
}
