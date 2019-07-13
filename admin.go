package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Squwid/bytegolf/database"
	"github.com/Squwid/bytegolf/questions"
)

// isAdmin checks to see if a user has admin status
func isAdmin(w http.ResponseWriter, req *http.Request) bool {
	if !database.InProd() {
		return true
	}
	if !loggedIn(w, req) {
		return false
	}
	// TODO: admin doesnt work with new users, need to fix using databases
	return false
}

func adminholes(w http.ResponseWriter, req *http.Request) {
	if !isAdmin(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	qs, err := questions.GetAllQuestions()
	if err != nil {
		logger.Printf("error getting all questions: %v\n", err)
		w.Write([]byte(err.Error()))
		return
	}

	tpl.ExecuteTemplate(w, "adminqs.html", struct {
		Questions       []questions.Question
		QuestionsAmount int
	}{
		Questions:       qs,
		QuestionsAmount: len(qs),
	})
}

func createQuestion(w http.ResponseWriter, req *http.Request) {
	if !isAdmin(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	if req.Method == http.MethodPost {
		logger.Printf("adding new question request\n")
		var name, question, answer, difficulty, source, input string
		var live bool
		name = req.FormValue("name")
		question = req.FormValue("question")
		input = req.FormValue("input")
		answer = req.FormValue("answer")
		difficulty = req.FormValue("difficulty")
		live = req.FormValue("live") == "true"
		source = req.FormValue("source")

		q := questions.NewQuestion(name, question, input, answer, difficulty, source, live)
		err := q.Store()
		if err != nil {
			logger.Println("error storing q:", err)
			w.Write([]byte(err.Error()))
			return
		}
	}
	http.Redirect(w, req, "/admin/holes", http.StatusSeeOther)
}

func deployQuestion(w http.ResponseWriter, req *http.Request) {
	if !isAdmin(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	if req.Method == http.MethodGet {
		s := strings.Split(req.URL.Path, "/")
		path := s[3]

		questions.MakeLive(path)
	}
	http.Redirect(w, req, "/admin/holes", http.StatusSeeOther)
}

func archiveQuestion(w http.ResponseWriter, req *http.Request) {
	if !isAdmin(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	if req.Method == http.MethodGet {
		s := strings.Split(req.URL.Path, "/")
		path := s[3]

		questions.ArchiveQuestion(path)
	}
	http.Redirect(w, req, "/admin/holes", http.StatusSeeOther)
}

func deletehole(w http.ResponseWriter, req *http.Request) {
	s := strings.Split(req.URL.Path, "/")
	path := s[3]
	err := questions.RemoveQuestion(path)
	if err != nil {
		logger.Printf("error deleting question %s: %v\n", path, err)
		w.Write([]byte(fmt.Sprintf("error deleting question %s: %v", path, err)))
		return
	}
	http.Redirect(w, req, "/admin/holes", http.StatusSeeOther)
	return
}

func admin(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/admin/holes", http.StatusSeeOther)
}
