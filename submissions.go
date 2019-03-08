package main

import (
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/Squwid/bytegolf/runner"
)

var lastSubmission = make(map[string]*runner.CodeResponse)
var lastHole = make(map[string]int)

func submission(w http.ResponseWriter, req *http.Request) {
	// action=""
	if req.Method != http.MethodPost {
		// TODO: handle here
		return
	}
	// declare the internal server error handling to show the user the error the occurred
	intErr := func() { http.Error(w, "an internal server error occurred", http.StatusInternalServerError) }

	// Get the question from the hole
	hole := strings.TrimLeft(req.URL.Path, "/submit/")
	/*
		// What should i do with this code?
		question, err := getHoleByLink(hole)
		if err != nil {
			http.Redirect(w, req, "/holes", http.StatusSeeOther)
			return
		}
	*/

	user, err := FetchUser(w, req)
	if err != nil {
		logger.Printf("error fetching user: %v\n", err)
		intErr()
		return
	}
	if !loggedIn(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	file, fileHead, err := req.FormFile("codefile")
	if err != nil {
		logger.Printf("an error occurred opening a form file: %v\n", err)
		intErr()
		return
	}
	defer file.Close()

	lang := req.FormValue("language") // todo: handle false languages in the future
	logger.Printf("submission %s %s %s\n", user, hole, lang)

	bs, err := ioutil.ReadAll(file) // buffer in the future?
	if err != nil {
		logger.Printf("error reading all file : %v\n", err.Error())
		intErr()
		return
	}

	// run the code from the input through the submission system
	runnerClient := runner.NewClient()
	submission := runner.NewCodeSubmission(user.Email, hole, fileHead.Filename, lang, string(bs), runnerClient, awsSess)
	runnerResp, err := submission.Send(true) // true stands for save local
	if err != nil {
		logger.Printf("error using code runner : %s\n", err.Error())
		intErr()
		return
	}

	lastSubmission[user.Email] = runnerResp
	http.Redirect(w, req, path.Join("play", hole), http.StatusSeeOther)
	return

	// This code was moved keep here until i know what to do with it
	/*
		score := question.Score(submission)

		// todo: make a translator function to take of this for me
		lbScore := LBSingleScore{
			Username: user.Email,
			Language: submission.Language,
			Score:    int(score),
		}
		addScore(hole, lbScore)

		playPage.ShowCorrect = true
		playPage.CorrectMessage = fmt.Sprintf("%s\nis the correct output.", runnerResp.Output)
		playPage.CurrentScore = int(score)
		// exeTpl(true, false, runnerResp.Output+"\n was the correct output.", "BEST SCORE HOLDER")
		exeTpl()
		return
	*/
}

// TODO: RENABLE THIS
func play(w http.ResponseWriter, req *http.Request) {
}

// 	if req.Method != http.MethodGet {
// 		w.Write([]byte("this method is not supported"))
// 		return
// 	}
// 	// Get the users current hole from the url
// 	var hole int
// 	holeStr := strings.TrimLeft(req.URL.Path, "/play/")
// 	if len(holeStr) == 0 {
// 		hole = 1
// 	} else {
// 		i, err := strconv.Atoi(holeStr)
// 		if err != nil {
// 			hole = 1
// 		} else {
// 			if i < 1 {
// 				hole = 1
// 			} else if i > 9 {
// 				hole = 9
// 			} else {
// 				hole = i
// 			}
// 		}
// 	}

// 	// error functions that are only needed in this scope
// 	intErr := func() { http.Error(w, "an internal server error occurred", http.StatusInternalServerError) }
// 	// playTpl holds the data for the play page, only needed in this scope
// 	type playTpl struct {
// 		Question          *questions.Question
// 		ShowNeverAnswered bool
// 		ShowIncorrect     bool
// 		IncorrectMessage  string
// 		ShowCorrect       bool
// 		CorrectMessage    string
// 		CurrentScore      int

// 		// Leaderboards
// 		FirstPlace  LBSingleScore
// 		SecondPlace LBSingleScore
// 		ThirdPlace  LBSingleScore
// 	}

// 	var playPage playTpl
// 	exeTpl := func() {
// 		first, second, third := getTopThree(hole)
// 		if first.Score != 0 {
// 			playPage.FirstPlace = *first
// 		}
// 		if second.Score != 0 {
// 			playPage.SecondPlace = *second
// 		}
// 		if third.Score != 0 {
// 			playPage.ThirdPlace = *third
// 		}
// 		tpl.ExecuteTemplate(w, "play.html", playPage)
// 	}

// 	question, err := getHoleByLink(hole)
// 	if err != nil {
// 		http.Redirect(w, req, "/holes", http.StatusSeeOther)
// 		return
// 	}

// 	playPage.Question = question

// 	if !loggedIn(w, req) {
// 		playPage.ShowNeverAnswered = true
// 	} else {
// 		user, err := FetchUser(w, req)
// 		if err != nil {
// 			logger.Printf("error fetching user: %v\n", err)
// 			intErr()
// 			return
// 		}
// 		_, idx, exist := userHasSubmission(hole, user.Email)
// 		if exist {
// 			playPage.CorrectMessage = fmt.Sprintf("You have already answered this hole correctly in %v bytes!", holeScores[hole][idx].Score)
// 			playPage.ShowCorrect = true
// 			playPage.CurrentScore = holeScores[hole][idx].Score
// 		} else {
// 			playPage.ShowNeverAnswered = true
// 		}
// 	}

// 	exeTpl()
// 	return
// }
