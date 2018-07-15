package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Squwid/bytegolf/bgaws"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// GolfResponse TODO
type golfResponse struct {
	User     *bgaws.User
	Name     string
	LoggedIn bool
	Game     *bgaws.Game
	GameName string
}

func index(w http.ResponseWriter, req *http.Request) {
	u := getUser(w, req)
	if req.Method == http.MethodPost {
		if !currentlyLoggedIn(w, req) {
			http.Redirect(w, req, "/login", http.StatusSeeOther)
			return
		}
		// reqHoles = req.FormValue("holes")
		// reqMaxPlayers = req.FormValue("maxplayers")
		// reqName = req.FormValue("gamename")
		// reqPass = req.FormValue("password")
		newGame(w, req)
		cookie := &http.Cookie{
			Name:  "gameid",
			Value: currentGame.ID,
		}
		http.SetCookie(w, cookie)

		http.Redirect(w, req, "/currentgame", http.StatusSeeOther)
		return
	}
	// if they send the get method
	tpl.ExecuteTemplate(w, "index.html", golfResponse{
		User:     u,
		Name:     u.Username,
		Game:     currentGame,
		LoggedIn: currentlyLoggedIn(w, req),
	})
}

func current(w http.ResponseWriter, req *http.Request) {
	cookie, _ := req.Cookie("gameid")
	fmt.Println("CurrentGameID:", currentGame.ID, "CookieString:", cookie.Value)
	// if the player is not currently logged in
	if !currentlyLoggedIn(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	// if the user is not in a game yet
	if !userInGame(w, req) {
		fmt.Println("user not in game")
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	user := getUser(w, req)
	tpl.ExecuteTemplate(w, "currentgame.html", golfResponse{
		User:     user,
		Name:     user.Username,
		Game:     currentGame,
		GameName: currentGame.ID,
		LoggedIn: currentlyLoggedIn(w, req),
	})
}

func leaderboard(w http.ResponseWriter, req *http.Request) {
	user := getUser(w, req)
	tpl.ExecuteTemplate(w, "leaderboards.html", golfResponse{
		User:     user,
		Name:     user.Username,
		LoggedIn: currentlyLoggedIn(w, req),
		Game:     currentGame,
	})
}

func dev(w http.ResponseWriter, req *http.Request) {
	user := getUser(w, req)
	if !currentlyLoggedIn(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	if user.Role != "dev" {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
}

func signup(w http.ResponseWriter, req *http.Request) {
	if currentlyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		un := req.FormValue("username")
		p := req.FormValue("password")
		// username taken?
		if bgaws.UserExist(un) {
			http.Error(w, "Username already taken", http.StatusForbidden)
			return
		}
		// create session
		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		// c.MaxAge = sessionLength
		http.SetCookie(w, c)
		currentSessions[c.Value] = session{un, time.Now()}
		// store user in dbUsers
		bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		u := &bgaws.User{
			Username: un,
			Password: string(bs),
		}
		bgaws.CreateUser(u)
		// redirect
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "signup.html", golfResponse{
		User:     nil,
		Name:     "",
		LoggedIn: true,
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
			http.Error(w, "That user does not exist", http.StatusForbidden)
			return
		}
		//TODO: encryption
		// err := bcrypt.CompareHashAndPassword(u.Password, p)
		user, _ := bgaws.GetUser(reqName)
		if user.Password != reqPass {
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
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	// if the user uses a get method just return the login screen with
	// no name on it
	tpl.ExecuteTemplate(w, "login.html", golfResponse{
		User:     nil,
		Game:     currentGame,
		LoggedIn: false,
	})
}
