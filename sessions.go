package main

import (
	"net/http"
	"time"

	"github.com/Squwid/bytegolf/bgaws"
	uuid "github.com/satori/go.uuid"
)

func newGame(w http.ResponseWriter, req *http.Request) {
	gameID, _ := uuid.NewV4()
	currentGame = &bgaws.Game{
		ID:          gameID.String(),
		Name:        req.FormValue("gamename"),
		StartedTime: time.Now(),
		Started:     true,
	}
	user := getUser(w, req)
	logger.Printf("%s created new game %s\n", user.Username, gameID.String())
}

func joinCurrentGame(w http.ResponseWriter, req *http.Request) {
	if !currentGame.Started {
		http.Error(w, "There are no current games", http.StatusForbidden)
		return
	}
	if !currentlyLoggedIn(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "gameid",
		Value: currentGame.ID,
	})
	user := getUser(w, req)
	logger.Printf("%s joined game %s\n", user.Username, currentGame.ID)
}

func userInGame(w http.ResponseWriter, req *http.Request) bool {
	cookie, err := req.Cookie("gameid")
	if err != nil {
		return false
	}
	// fmt.Println("CurrentGameID:", currentGame.ID, "CookieString:", cookie.String())
	if currentGame.ID != cookie.Value {
		return false
	}
	return true
}

func getUser(w http.ResponseWriter, req *http.Request) *bgaws.User {
	var cookie *http.Cookie
	cookie, err := req.Cookie("session")
	if err != nil {
		session, _ := uuid.NewV4()
		cookie = &http.Cookie{
			Name:  "session",
			Value: session.String(),
		}
	}
	// http.SetCookie(w, cookie)

	// if the user exists already, get user
	u, _ := bgaws.GetUser(currentSessions[cookie.Value].Username) //TODO: when does this error get called?
	return u
}

// currentlyLoggedIn checks to see if a user is currently logged in already by
// checking their session ID
func currentlyLoggedIn(w http.ResponseWriter, req *http.Request) bool {
	cookie, err := req.Cookie("session")
	if err != nil {
		return false
	}
	session, ok := currentSessions[cookie.Value]
	if ok {
		// if session currently exists
		session.lastActivity = time.Now()
		currentSessions[cookie.Value] = session // todo: does this do what i think it does?
	}
	// refresh session
	// c.MaxAge = sessionLength
	// http.SetCookie(w, cookie)
	return ok
}
