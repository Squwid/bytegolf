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
	tpl.ExecuteTemplate(w, "account.html", struct{}{})
}

func leaderboards(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "leaderboards.html", struct{}{})
}

func login(w http.ResponseWriter, req *http.Request) {
	exeTpl := func(incorrectPass bool) {
		tpl.ExecuteTemplate(w, "login.html", struct {
			IncorrectPassword bool
		}{
			IncorrectPassword: incorrectPass,
		})
	}

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
			// incorrect password == true
			exeTpl(true)
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
	// incorrect password == false
	exeTpl(false)
}

func holes(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "holes.html", struct {
		Questions map[int]questions.Question
	}{
		Questions: qs,
	})
}
