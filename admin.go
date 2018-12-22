package main

import (
	"net/http"
	"strings"

	"github.com/Squwid/bytegolf/aws"

	"github.com/Squwid/bytegolf/questions"
)

func admin(w http.ResponseWriter, req *http.Request) {
	// if !loggedIn(w, req) {
	// 	http.Redirect(w, req, "/login", http.StatusSeeOther)
	// 	return
	// }
	// todo: check to see if the user is an admin and allowed to see the page
	if req.Method == http.MethodPost {
		q := questions.NewQuestion(req.FormValue("name"), req.FormValue("question"), req.FormValue("answer"), req.FormValue("difficulty"), req.FormValue("source"), req.FormValue("link"))
		q.Store(true)
		qs = questions.ToMap(questions.GetLocalQuestions())
		logger.Printf("creating question %s\n", q.Name)
		http.Redirect(w, req, "/admin/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "admin.html", struct {
		Questions       map[int]questions.Question
		QuestionsAmount int
	}{
		Questions:       qs,
		QuestionsAmount: len(qs),
	})
}

func createuser(w http.ResponseWriter, req *http.Request) {
	logger.Printf("adding user request\n")
	if req.Method == http.MethodPost {
		email := req.FormValue("email")
		username := req.FormValue("username")
		password := req.FormValue("password")
		role := req.FormValue("role")
		_ = role // get rid of role as its not used right now
		user := aws.NewUser(email, username, password)
		err := user.Store()
		if err != nil {
			logger.Printf("error adding user %s to aws: %v\n", username, err)
		} else {
			logger.Printf("successfully added %s to aws\n", username)
		}
	}
	http.Redirect(w, req, "/admin/", http.StatusSeeOther)
}

func deletehole(w http.ResponseWriter, req *http.Request) {
	s := strings.Split(req.URL.Path, "/")
	path := s[3]
	RemoveQuestion(path)
	qs = questions.ToMap(questions.GetLocalQuestions())
	http.Redirect(w, req, "/admin/", http.StatusSeeOther)
	return
}

// RemoveQuestion removes a specific quetsion using a link to the question
func RemoveQuestion(link string) {
	logger.Printf("trying to remove question by the link of %s\n", link)
	for _, q := range qs {
		if q.Link == link {
			err := q.Remove()
			if err != nil {
				panic(err)
			}
		}
	}
}

func removeGame(gameID string) error {
	return nil
}
