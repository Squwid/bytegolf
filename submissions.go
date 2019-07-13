package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Squwid/bytegolf/questions"
	"github.com/Squwid/bytegolf/runner"
)

func submission(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Redirect(w, req, path.Join("play"), http.StatusSeeOther)
		return
	}

	// Do user stuff before the logic
	if !loggedIn(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	// declare the internal server error handling to show the user the error the occurred
	intErr := func() { http.Error(w, "an internal server error occurred", http.StatusInternalServerError) }

	user, err := fetchUser(w, req)
	if err != nil {
		logger.Printf("error fetching user: %v\n", err)
		intErr()
		return
	}

	// POST "/submit/{holeid}"
	holeID := strings.Split(req.URL.Path, "/")[2]
	logger.Printf("holeID: %v\tSplit: %v\n", holeID, strings.Split(req.URL.Path, "/"))

	question, err := questions.GetQuestionByID(holeID)
	if err != nil {
		logger.Printf("error trying to get question %v\n", holeID)
		http.Redirect(w, req, "/holes/", http.StatusSeeOther)
		return
	}
	if !question.Live {
		logger.Printf("tried to get %v hole that isnt live\n", holeID)
		http.Redirect(w, req, "/holes/", http.StatusSeeOther)
		return
	}

	// parse the uploaded file
	file, fileHead, err := req.FormFile("codefile")
	if err != nil {
		logger.Printf("an error occurred opening a form file: %v\n", err)
		intErr()
		return
	}
	defer file.Close()

	lang := req.FormValue("language") // todo: handle false languages in the future
	logger.Printf("new submission %v %v (Lang): %s\n", user.ID, question.Name, lang)

	// read the file
	bs, err := ioutil.ReadAll(file) // buffer in the future?
	if err != nil {
		logger.Printf("error reading all file : %v\n", err.Error())
		intErr()
		return
	}

	// run the code from the input through the submission system
	// TODO: pass entire question through the runner rather than just the id to not have to make multiple calls to get the question (not a big deal for now just double logs)
	runnerClient := runner.NewClient() // todo: New way to do a runner rather than each call
	submission := runner.NewCodeSubmission(strconv.Itoa(user.ID), user.Username, question.ID, question.Input, fileHead.Filename, lang, string(bs), runnerClient, awsSess)
	_, err = submission.Send(true) // true stands for save local
	if err != nil {
		logger.Printf("error using code runner : %s\n", err.Error())
		intErr()
		return
	}
	time.Sleep(1 * time.Second) // let the page update
	http.Redirect(w, req, "/play/"+question.ID, http.StatusSeeOther)
	return
}

// play shows a specific hole, using a number at the end of the request. It then will add it as the last hole to the map, using that as the post request as what hole the player submitted
func play(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.Write([]byte("this method is not supported"))
		return
	}
	// internal server error gets called often
	intErr := func() { http.Error(w, "an internal server error occurred", http.StatusInternalServerError) }

	type player struct {
		Place    string
		Username string
		Language string
		Score    int
	}

	type playTpl struct {
		Question          *questions.Question
		ShowNeverAnswered bool
		ShowIncorrect     bool
		IncorrectMessage  string
		IncorrectSub      string
		ShowCorrect       bool
		CorrectMessage    string
		CurrentScore      int
		Player            player

		// Leaderboards
		Leaderboard map[int]runner.LbOverview
	}

	// assign the hole from the path
	holeID := strings.Split(req.URL.Path, "/")[2]

	var playPage playTpl
	exeTpl := func() {
		playPage.Leaderboard = runner.GetHoleLB(holeID)
		tpl.ExecuteTemplate(w, "play.html", playPage)
	}

	question, err := questions.GetQuestionByID(holeID)
	if err != nil {
		logger.Printf("error trying to get question %v\n", holeID)
		http.Redirect(w, req, "/holes/", http.StatusSeeOther)
		return
	}
	if !question.Live {
		logger.Printf("tried to get %v hole that isnt live\n", holeID)
		http.Redirect(w, req, "/holes/", http.StatusSeeOther)
		return
	}

	playPage.Question = question

	if loggedIn(w, req) {
		// the player is logged in so grab them
		user, err := fetchUser(w, req)
		if err != nil {
			logger.Printf("error fetching user: %v\n", err)
			intErr()
			return
		}
		prev := runner.PreviouslyAnswered(holeID, strconv.Itoa(user.ID))
		if prev.Correct {
			// TODO: add a date to the correct screen
			playPage.ShowCorrect = true
			playPage.CorrectMessage = fmt.Sprintf("You have gotten this answer correct using %s in %v bytes!", prev.Language, prev.Score)
			playPage.CurrentScore = prev.Score

			// for the custom leaderboard spot
			playPage.Player.Language = prev.Language
			playPage.Player.Username = user.Username
			playPage.Player.Score = prev.Score
			playPage.Player.Place = "--" // TODO: make a current place
		} else {
			playPage.ShowNeverAnswered = true
		}
		if runner.LS.IsIncorrect(user.Username, holeID) {
			playPage.ShowIncorrect = true
			playPage.IncorrectMessage = fmt.Sprintf("Your last submission of this hole at %s is incorrect", runner.LS.GetTime(user.Username).Format("Jan 2 15:04 MST 2006"))
			playPage.IncorrectSub = runner.LS.GetOutput(user.Username)
		}
	} else {
		// not logged in so show nothing
		playPage.ShowNeverAnswered = false
	}
	exeTpl()
	return
}
