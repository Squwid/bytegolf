package main

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Squwid/bytegolf/aws"
	"github.com/Squwid/bytegolf/runner"
)

// MessageOfTheDay returns a string that contains message of the day on the index.html page
func MessageOfTheDay() string {
	var messages = []string{"Bugs not birds..."}
	// random := rand.Intn(len(messages) - 1)
	return messages[0]
}

func index(w http.ResponseWriter, req *http.Request) {
	logger.Printf("here")
	logger.Println(tpl.ExecuteTemplate(w, "index.html", struct {
		MOTD string
	}{
		MOTD: MessageOfTheDay(),
	}))
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

func play(w http.ResponseWriter, req *http.Request) {
	internalErr := func(w http.ResponseWriter) {
		http.Error(w, "an internal server error occurred", http.StatusInternalServerError)
	}
	userErr := func(w http.ResponseWriter) {
		http.Error(w, "user misuse error", http.StatusBadRequest)
	}
	if !loggedIn(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	user, err := FetchUser(w, req)
	if err != nil {
		logger.Printf("error fetching user: %v\n", err)
		return
	}
	logger.Printf("user fetched: %v\n", user)
	if !CurrentGame.InProgress() {
		http.Redirect(w, req, "/create", http.StatusSeeOther)
		return
	}

	hole := getUserHole(req)
	question := CurrentGame.Questions[hole]

	// if the player has submitted a bytegolf submission
	if req.Method == "POST" {
		lang := req.FormValue("language")
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
		// TODO: Check output file and scoring system
		if !CurrentGame.CheckSubmission(hole, runnerResp.Output) {
			// the users submission is wrong
			//todo: handle this
		}
		_ = submission
		_ = runnerResp
	}

	tpl.ExecuteTemplate(w, "play.html", struct {
		Game     Game
		User     aws.User
		Hole     int
		Question aws.Question

		HolesCorrect int
		TotalScore   int
	}{
		Game:         CurrentGame,
		User:         user,
		HolesCorrect: 3,
		TotalScore:   250,
		Hole:         hole,
		Question:     question,
	})
}

func create(w http.ResponseWriter, req *http.Request) {
	br := func() {
		http.Error(w, "bad request", http.StatusBadRequest)
	}
	if !loggedIn(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	if CurrentGame.InProgress() {
		http.Redirect(w, req, "/play", http.StatusSeeOther)
		return
	}
	if req.Method == "POST" {
		// if the user is trying to create a game
		holes, err := strconv.Atoi(req.FormValue("holes"))
		if err != nil {
			br()
			return
		}
		mp, err := strconv.Atoi(req.FormValue("maxplayers"))
		if err != nil {
			br()
			return
		}
		g, err := NewGame(req.FormValue("gamename"), req.FormValue("password"), "medium", holes, mp)
		if err != nil {
			panic(err)
		}
		CurrentGame = g
		if CurrentGame.Start() != nil {
			panic(err)
		}
		logger.Printf("new game was successfully created\n")
		http.Redirect(w, req, "/play", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "currentgames.html", struct{}{})
}
