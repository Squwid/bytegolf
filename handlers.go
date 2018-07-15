package main

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Squwid/bytegolf/bgaws"
	"github.com/Squwid/bytegolf/runner"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func index(w http.ResponseWriter, req *http.Request) {
	user := getUser(w, req)
	// If there is a post on index that means that the user is
	// creating a new game
	if req.Method == http.MethodPost {
		if !currentlyLoggedIn(w, req) {
			http.Redirect(w, req, "/login", http.StatusSeeOther)
			return
		}
		err := CreateNewGame(w, req)
		if err != nil {
			logger.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		gameCookie := &http.Cookie{
			Name:  "gameid",
			Value: currentGame.ID,
		}
		http.SetCookie(w, gameCookie)
		http.Redirect(w, req, "/currentgame", http.StatusSeeOther)
		return
	}
	// if they send the get method
	tpl.ExecuteTemplate(w, "index.html", golfResponse{
		User:     user,
		Name:     user.Username,
		Game:     currentGame,
		LoggedIn: currentlyLoggedIn(w, req),
	})
}

func current(w http.ResponseWriter, req *http.Request) {
	// if the player isnt logged in send them to the login screen
	if !currentlyLoggedIn(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	user := getUser(w, req)

	if !userInGame(w, req) {
		if !currentGame.Started {
			http.Error(w, "there is not a current game", http.StatusNoContent)
			// if the game exists
			gameCookie := &http.Cookie{
				Name:  "gameid",
				Value: currentGame.ID,
			}
			http.SetCookie(w, gameCookie)
		}
		gameCookie := &http.Cookie{
			Name:  "gameid",
			Value: currentGame.ID,
		}
		currentGame.CurrentPlayers++
		logger.Printf("%s added to game %s\n", user.Username, currentGame.Name)
		logger.Printf("there are now %v people in game %s\n", currentGame.CurrentPlayers, currentGame.Name)
		http.SetCookie(w, gameCookie)
	}

	if req.Method == http.MethodPost {
		// open submitted file
		file, fileHead, err := req.FormFile("codefile")
		lang := req.FormValue("language")
		if err != nil {
			logger.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		logger.Printf("%s uploading file %s\n", user.Username, fileHead.Filename)

		// read
		bs, err := ioutil.ReadAll(file)
		if err != nil {
			logger.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		client := runner.NewClient()
		sub := runner.NewCodeSubmission(lang, string(bs), client)
		_, err = sub.Send() // todo: reply goes here but we dont deal with it yet
		if err != nil {
			logger.Println(err.Error())
			http.Error(w, "an unexpected error has occured", http.StatusInternalServerError)
			return
		}
	}
	_, err := req.Cookie("gameid")
	if err != nil {
		logger.Println(err.Error())
		http.Error(w, "an unexpected error has occured", http.StatusInternalServerError)
		return
	}
	// fmt.Println("CurrentGameID:", currentGame.ID, "CookieString:", cookie.Value)
	tpl.ExecuteTemplate(w, "currentgame.html", golfResponse{
		User:     user,
		Name:     user.Username,
		Game:     currentGame,
		GameName: currentGame.Name,
		LoggedIn: currentlyLoggedIn(w, req),
	})
}

func leaderboard(w http.ResponseWriter, req *http.Request) {
	// if the player isnt logged in send them to the login screen
	if !currentlyLoggedIn(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	user := getUser(w, req)
	tpl.ExecuteTemplate(w, "leaderboards.html", golfResponse{
		User:     user,
		Name:     user.Username,
		LoggedIn: currentlyLoggedIn(w, req),
		Game:     currentGame,
	})
}

func dev(w http.ResponseWriter, req *http.Request) {
	// if the player isnt logged in send them to the login screen
	if !currentlyLoggedIn(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	user := getUser(w, req)
	logger.Printf("%s is trying to access DEV page\n", user.Username)
	if user.Role != "dev" {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "devtools.html", golfResponse{
		User:     user,
		Name:     user.Username,
		LoggedIn: currentlyLoggedIn(w, req),
		Game:     currentGame,
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
		// todo: encrypt passwords

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
		Game:     currentGame,
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
		// if user.Password != reqPass {
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
		Game:     currentGame,
		LoggedIn: false,
	})
}
