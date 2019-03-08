package main

import (
	"net/http"
	"strings"

	"github.com/Squwid/bytegolf/questions"
)

func submission(w http.ResponseWriter, req *http.Request) {
	// if req.Method != http.MethodPost {
	// 	http.Redirect(w, req, path.Join("play"), http.StatusSeeOther)
	// 	return
	// }
	// // declare the internal server error handling to show the user the error the occurred
	// intErr := func() { http.Error(w, "an internal server error occurred", http.StatusInternalServerError) }

	// if !loggedIn(w, req) {
	// 	http.Redirect(w, req, "/login", http.StatusSeeOther)
	// 	return
	// }
	// // hole := getUserHole(req) // get the hole from the cookie

	// user, err := FetchUser(w, req)
	// if err != nil {
	// 	logger.Printf("error fetching user: %v\n", err)
	// 	intErr()
	// 	return
	// }

	// file, fileHead, err := req.FormFile("codefile")
	// if err != nil {
	// 	logger.Printf("an error occurred opening a form file: %v\n", err)
	// 	intErr()
	// 	return
	// }
	// defer file.Close()

	// lang := req.FormValue("language") // todo: handle false languages in the future
	// logger.Printf("submission %s %v %s\n", user, hole, lang)

	// bs, err := ioutil.ReadAll(file) // buffer in the future?
	// if err != nil {
	// 	logger.Printf("error reading all file : %v\n", err.Error())
	// 	intErr()
	// 	return
	// }

	// // run the code from the input through the submission system
	// runnerClient := runner.NewClient()
	// submission := runner.NewCodeSubmission(user.Email, strconv.Itoa(hole), fileHead.Filename, lang, string(bs), runnerClient, awsSess)
	// runnerResp, err := submission.Send(true) // true stands for save local
	// if err != nil {
	// 	logger.Printf("error using code runner : %s\n", err.Error())
	// 	intErr()
	// 	return
	// }
	// _ = runnerResp

	// http.Redirect(w, req, path.Join("play", strconv.Itoa(hole)), http.StatusSeeOther)
	// return

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

// play shows a specific hole, using a number at the end of the request. It then will add it as the last hole to the map, using that as the post request as what hole the player submitted
func play(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.Write([]byte("this method is not supported"))
		return
	}
	// internal server error gets called often
	intErr := func() { http.Error(w, "an internal server error occurred", http.StatusInternalServerError) }
	// playTpl holds the data for the play page, only needed in this scope
	type LBSingleScore struct {
		Username string
		Language string
		Score    string
	}
	type playTpl struct {
		Question          *questions.Question
		ShowNeverAnswered bool
		ShowIncorrect     bool
		IncorrectMessage  string
		ShowCorrect       bool
		CorrectMessage    string
		CurrentScore      int

		// Leaderboards
		FirstPlace  LBSingleScore
		SecondPlace LBSingleScore
		ThirdPlace  LBSingleScore
	}

	// assign the hole from the path
	holeID := strings.Split(req.URL.Path, "/")[2]

	var playPage playTpl
	exeTpl := func() {
		// TODO: Leaderboards
		tpl.ExecuteTemplate(w, "play.html", playPage)
	}

	question := questions.GetQuestion(holeID)
	if question == nil {
		logger.Printf("tried to get %s hole that isnt live\n", holeID)
		http.Redirect(w, req, "/holes/", http.StatusSeeOther)
		return
	}

	playPage.Question = question

	if loggedIn(w, req) {
		// the player is logged in so grab them
		user, err := FetchUser(w, req)
		if err != nil {
			logger.Printf("error fetching user: %v\n", err)
			intErr()
			return
		}
		playPage.ShowNeverAnswered = true
		_ = user
	} else {
		// not logged in so show nothing
		playPage.ShowNeverAnswered = false
	}
	exeTpl()
	return
}

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
