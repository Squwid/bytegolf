package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Squwid/bytegolf/bgaws"
	uuid "github.com/satori/go.uuid"
)

// CreateNewGame todo
func CreateNewGame(w http.ResponseWriter, req *http.Request) error {
	gameID, _ := uuid.NewV4()
	holes, err := strconv.Atoi(req.FormValue("holes"))
	if err != nil {
		logger.Println(err)
		return err
	}
	holes = 3 //todo: remove this. It is hardset to 3 right now because there are not enough
	// questions in the bank
	max, err := strconv.Atoi(req.FormValue("maxplayers"))
	if err != nil {
		logger.Println(err)
		return err
	}
	diff := req.FormValue("difficulty")
	logger.Printf("new game requested with %v holes at %s difficulty\n", holes, diff)
	diff = "medium" // todo: this is the only difficulty that we currently have
	qs, err := bgaws.GetQuestions(diff, 3)
	if err != nil {
		return err
	}
	currentGame = Game{
		ID:             gameID.String(),
		Name:           req.FormValue("gamename"),
		Password:       req.FormValue("password"),
		CurrentPlayers: 0,
		MaxPlayers:     max,
		Holes:          holes,
		Difficulty:     diff,
		StartedTime:    time.Now(),
		Started:        true,
		Questions:      qs,
	}
	user := getUser(w, req)
	logger.Printf("%s created new game %s\n", user.Username, currentGame.Name)
	return nil
}

// userInGame returns a bool whether or not the player sending the request
// is currently in a game
func userInGame(w http.ResponseWriter, req *http.Request) bool {
	cookie, err := req.Cookie("gameid")
	if err != nil {
		return false
	}
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
