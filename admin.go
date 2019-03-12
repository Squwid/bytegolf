package main

import (
	"net/http"
	"strings"

	"github.com/Squwid/bytegolf/aws"
	"github.com/Squwid/bytegolf/questions"
)

// isAdmin checks to see if a user has admin status
func isAdmin(w http.ResponseWriter, req *http.Request) bool {
	if !loggedIn(w, req) {
		return false
	}
	user, err := FetchUser(w, req)
	if err != nil {
		return false
	}
	if user.Role != aws.RoleAdmin && user.Role != aws.RoleDev {
		return false
	}
	return true
}

func adminholes(w http.ResponseWriter, req *http.Request) {
	if !isAdmin(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	qs, err := questions.GetLocalQuestions()
	if err != nil {
		logger.Printf("error getting local questions: %v\n", err)
		w.Write([]byte(err.Error()))
		return
	}
	var qsMap = map[int]questions.Question{}
	for i, q := range qs {
		qsMap[i+1] = q
	}
	tpl.ExecuteTemplate(w, "adminqs.html", struct {
		Questions       map[int]questions.Question
		QuestionsAmount int
	}{
		Questions:       qsMap,
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

		allQs, err := questions.GetLocalQuestions()
		if err != nil {
			logger.Println("err getting all local qs", err)
			w.Write([]byte(err.Error()))
			return
		}
		for i, q := range allQs {
			if q.ID == path {
				err = allQs[i].Deploy()
				if err != nil {
					logger.Println("error deploying q:", err)
					w.Write([]byte(err.Error()))
				}
				break
			}
		}
	}
	http.Redirect(w, req, "/admin/holes", http.StatusSeeOther)
}

func refreshQuestions(w http.ResponseWriter, req *http.Request) {
	if !isAdmin(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	questions.UpdateQuestions()
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

		allQs, err := questions.GetLocalQuestions()
		if err != nil {
			logger.Println("err getting all local qs", err)
			w.Write([]byte(err.Error()))
			return
		}
		for i, q := range allQs {
			if q.ID == path {
				err = allQs[i].RemoveLive()
				if err != nil {
					logger.Println("error storing q:", err)
					w.Write([]byte(err.Error()))
				}
				break
			}
		}
	}
	http.Redirect(w, req, "/admin/holes", http.StatusSeeOther)
}

func createUser(w http.ResponseWriter, req *http.Request) {
	if !isAdmin(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	if req.Method == http.MethodPost {
		logger.Printf("adding user request\n")
		var email, username, password, role string
		email = req.FormValue("email")
		username = req.FormValue("username")
		password = req.FormValue("password")
		role = req.FormValue("role")

		user := aws.NewUser(email, username, role, password)
		err := user.Store()
		if err != nil {
			logger.Printf("error adding user %s to aws: %v\n", username, err)
		} else {
			logger.Printf("successfully added %s to aws\n", username)
		}
	}
	http.Redirect(w, req, "/admin/holes", http.StatusSeeOther)
}

// func deletehole(w http.ResponseWriter, req *http.Request) {
// 	s := strings.Split(req.URL.Path, "/")
// 	path := s[3]
// 	RemoveQuestion(path)
// 	// qs = questions.ToMap(questions.GetLocalQuestions())
// 	http.Redirect(w, req, "/admin/holes", http.StatusSeeOther)
// 	return

func admin(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/admin/holes", http.StatusSeeOther)
}

// // RemoveQuestion removes a specific quetsion using a link to the question
// func RemoveQuestion(link string) {
// 	logger.Printf("trying to remove question by the link of %s\n", link)
// 	for _, q := range qs {
// 		if q.Link == link {
// 			// err := q.Remove()
// 			if err != nil {
// 				panic(err)
// 			}
// 		}
// 	}
// }

// func removeGame(gameID string) error {
// 	return nil
// }
