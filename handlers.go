package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Squwid/bytegolf/bgaws"
	"github.com/Squwid/bytegolf/runner"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

/* FIXME: Currently I am returning a structure to each part of the handler but
instead I should be returning an anon struct back */

func index(w http.ResponseWriter, req *http.Request) {
	user := getUser(w, req)
	// If there is a post on index that means that the user is
	// creating a new game
	if req.Method == http.MethodPost {
		if !currentlyLoggedIn(w, req) {
			http.Redirect(w, req, "/login", http.StatusSeeOther)
			return
		}
		// if there is already a current game
		if CurrentGame.Started {
			http.Redirect(w, req, "/currentgame", http.StatusSeeOther)
			return
		}
		err := CreateNewGame(w, req)
		if err != nil {
			logger.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = addUserToCurrent(w, *user)
		if err != nil {
			logger.Println(err.Error())
			http.Error(w, "an internal error occurred", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, req, "/currentgame", http.StatusSeeOther)
		return
	}
	// if they send the get method
	tpl.ExecuteTemplate(w, "index.html", golfResponse{
		User:     user,
		Name:     user.Username,
		Game:     CurrentGame,
		LoggedIn: currentlyLoggedIn(w, req),
	})
}

// current holds the code for the submission along with joining a current game
func current(w http.ResponseWriter, req *http.Request) {
	// if the player isnt logged in send them to the login screen
	if !currentlyLoggedIn(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	// the user is logged in so get the user
	user := getUser(w, req)

	// if the user is not in a game check to see if there is a current game
	if !userInGame(w, req) {
		if !CurrentGame.Started {
			http.Error(w, "there is not a current game", http.StatusNoContent)
			return
		}
		// if the user is not in a game and there is a current game add them to it
		err := addUserToCurrent(w, *user)
		if err != nil {
			http.Error(w, "an internal error occurred", http.StatusInternalServerError)
		}
	}

	// parse what hole the player is on
	var hole int
	h := strings.TrimPrefix(req.URL.Path, "/currentgame/")
	if len(h) == 1 {
		i, err := strconv.Atoi(h)
		hole = i
		if err != nil {
			hole = 1
		}
	} else {
		hole = 1
	}
	if hole > CurrentGame.Holes {
		hole = CurrentGame.Holes // if the user goes over the limit set the hole to the max
	}

	// if the user is submitting a file
	if req.Method == http.MethodPost {
		// open submitted file
		lang := req.FormValue("language")
		file, fileHead, err := req.FormFile("codefile")
		if err != nil {
			logger.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		logger.Printf("%s uploading file %s\n", user.Username, fileHead.Filename)

		// read the file
		bs, err := ioutil.ReadAll(file)
		if err != nil {
			logger.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Run the code given from the submission
		client := runner.NewClient()
		sub := runner.NewCodeSubmission(user.Username, CurrentGame.ID, fileHead.Filename, lang, string(bs), client)
		resp, err := sub.Send()
		if err != nil {
			logger.Println(err.Error())
			http.Error(w, "an unexpected error has occured", http.StatusInternalServerError)
			return
		}
		q, ok := CurrentGame.Questions[hole]
		if !ok {
			http.Error(w, "unable to find that hole", http.StatusInternalServerError)
		}
		correct := checkResponse(resp, &q)
		_ = correct //TODO: continue here
	}

	q, ok := CurrentGame.Questions[hole]
	if !ok {
		http.Error(w, "unable to find that hole", http.StatusInternalServerError)
		return
	}

	tpl.ExecuteTemplate(w, "currentgame.html", golfResponse{
		User:     user,
		Name:     user.Username,
		Game:     CurrentGame,
		GameName: CurrentGame.Name,
		Hole:     hole,
		Question: q,
		LoggedIn: currentlyLoggedIn(w, req),
	})
}

func signup(w http.ResponseWriter, req *http.Request) {
	// if the user is already logged in then send them to the home screen
	if currentlyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		reqName := req.FormValue("username")
		reqPassword := req.FormValue("password")

		// Check if username is already taken
		if bgaws.UserExist(reqName) {
			logger.Printf("user tried to register with %s but it was taken\n", reqName)
			http.Error(w, "Username already taken", http.StatusForbidden)
			return
		}
		// username is not taken so sign them up and create session
		bs, err := bcrypt.GenerateFromPassword([]byte(reqPassword), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		newUser := &bgaws.User{
			Username: reqName,
			// Password: reqPassword,
			Password: string(bs),
		}
		err = bgaws.CreateUser(newUser)
		if err != nil {
			logger.Println(err.Error())
			http.Error(w, "an unexpected error occured", http.StatusInternalServerError)
			return
		}

		sessionID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sessionID.String(),
		}

		http.SetCookie(w, c)
		currentSessions[c.Value] = session{
			Username:     newUser.Username,
			lastActivity: time.Now(),
		}
		logger.Printf("%s successfully signed up\n", newUser.Username)
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "signup.html", golfResponse{
		Game:     CurrentGame,
		LoggedIn: false,
	})
}

func login(w http.ResponseWriter, req *http.Request) {
	if currentlyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	// process form submission
	if req.Method == http.MethodPost {
		reqName := req.FormValue("username")
		reqPass := req.FormValue("password")

		if !bgaws.UserExist(reqName) {
			logger.Printf("user tried to login with %s but it does not exist\n", reqName)
			http.Error(w, "That user does not exist", http.StatusForbidden)
			return
		}

		user, _ := bgaws.GetUser(reqName)
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqPass))
		if err != nil {
			logger.Printf("%s tried to login with incorrect password\n", user.Username)
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}

		// Password entered was correct so create new session
		sessionID, _ := uuid.NewV4()
		cookie := &http.Cookie{
			Name:  "session",
			Value: sessionID.String(),
		}

		http.SetCookie(w, cookie)
		currentSessions[cookie.Value] = session{
			Username:     reqName,
			lastActivity: time.Now(),
		}
		logger.Printf("%s successfully logged in\n", user.Username)
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "login.html", golfResponse{
		Game:     CurrentGame,
		LoggedIn: false,
	})
}

func logout(w http.ResponseWriter, req *http.Request) {
	if !currentlyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	user := getUser(w, req)
	sessionCookie, err := req.Cookie("session")
	if err != nil {
		http.Error(w, "an unexpected error has occured", http.StatusInternalServerError)
		return
	}
	// todo: delete game cookie on logout
	delete(currentSessions, sessionCookie.Value)
	sessionCookie = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	logger.Printf("%s successfully logged out\n", user.Username)
	http.Redirect(w, req, "/login", http.StatusSeeOther)
}
