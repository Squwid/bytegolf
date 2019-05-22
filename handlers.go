package main

import (
	"net/http"

	"github.com/Squwid/bytegolf/questions"
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
	user, err := fetchUser(w, req)
	if err != nil {
		logger.Fatalln("error fetching user:", err)
		http.Error(w, "an internal server error has occurred", http.StatusInternalServerError)
		return
	}

	tpl.ExecuteTemplate(w, "profile.html", user)
}

func leaderboards(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "leaderboards.html", struct{}{})
}

func holes(w http.ResponseWriter, req *http.Request) {
	qs := questions.GetAllQuestions()
	var m = make(map[int]questions.Question)
	for k, v := range qs {
		m[k] = *v
	}
	tpl.ExecuteTemplate(w, "holes.html", struct {
		Questions map[int]questions.Question
	}{
		Questions: m,
	})
}
