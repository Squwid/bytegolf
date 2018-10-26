package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

// TutGames holds all of the tutorial games using their email as a key
var TutGames = make(map[string]*TutGame)

// TutHoles is all of the tutorial holes for the tutorial game mode
var TutHoles = make(map[int]*TutQuestion)

// TutQuestion is the entire question that goes into the hole struct
type TutQuestion struct {
	Question   string `json:"question"`
	Answer     string `json:"answer"`
	Difficulty string `json:"difficulty"`
}

// TutGame is a tutorial game that is used for the tutorial mode in bytegolf
type TutGame struct {
	TimeCreated time.Time
	Ended       bool
	TimeEnded   time.Time
	Scores      map[int]int
	Correct     map[int]bool
}

// InitTutorialQs creates the TutHoles map with the 3 questions that will be in the tutorial
func InitTutorialQs() {
	q1 := &TutQuestion{
		Question:   `In the shortest amount of code, print "Hello, World!" to console.`,
		Answer:     "Hello, World!",
		Difficulty: "easy",
	}
	q2 := &TutQuestion{
		Question:   `Find the total of every odd number from 1 to 1 million. Then print it to console.`,
		Answer:     "500000500000",
		Difficulty: "easy",
	}
	q3 := &TutQuestion{
		Question:   `By listing the first six prime numbers: 2, 3, 5, 7, 11, and 13, we can see that the 6th prime is 13. Write code to find and print the 10001st prime number.`,
		Answer:     "104743",
		Difficulty: "medium",
	}
	TutHoles[1] = q1
	TutHoles[2] = q2
	TutHoles[3] = q3
}

// NewTutorial is a new tutorial game for bytegolf users
func NewTutorial() *TutGame {
	var scores = make(map[int]int)
	scores[1] = 0
	scores[2] = 0
	scores[3] = 0
	return &TutGame{
		TimeCreated: time.Now(),
		Scores:      scores,
	}
}

// Add adds a game to the list of current games and it will stay there for 2 hours
func (game *TutGame) Add(name string) {
	TutGames[name] = game
	fmt.Printf("added game %s to games. now there are %d\n", name, len(TutGames))
}

// TutGameExist checks to see if a tutorial game for a specific user exists, and it returns the tutgame as well
func TutGameExist(name string) bool {
	if _, ok := TutGames[name]; ok {
		return true
	}
	return false
}

// GetTutGame gets a tutorial game from a specific username
func GetTutGame(name string) (*TutGame, error) {
	if game, ok := TutGames[name]; ok {
		return game, nil
	}
	return nil, errors.New("tutorial game does not exist for that player")
}

// InProgress checks to see if a tut game is currently in progress
func (game *TutGame) InProgress() bool {
	if game == nil {
		return false
	}
	if game.Ended {
		return false
	}
	return true
}

// tutJanitor loops every minute and deletes all games that currently exist in the map
func tutJanitor() {
	allocatedTime := 1 * time.Hour
	checkTime := 1 * time.Minute
	logger.Printf("starting tutorial game janitor checking every %v\n", checkTime)
	for range time.Tick(checkTime) {
		// fmt.Printf("checking now...\n")
		if len(TutGames) == 0 {
			continue
		}
		for k, v := range TutGames {
			// check to see if they are 1 hour long if so end the game
			if !v.Ended && time.Since(v.TimeCreated) >= allocatedTime {
				v.Ended = true // end the game
				v.TimeEnded = time.Now()
				fmt.Printf("stopping game %s.", k)
				continue
			}
			if v.Ended && time.Since(v.TimeEnded) >= allocatedTime {
				delete(TutGames, k) // delete the entry after anohter
				fmt.Printf("deleted game %s. now there are %d\n", k, len(TutGames))
				continue
			}
		}
	}
}

func tutCreator(w http.ResponseWriter, req *http.Request) {
	intErr := func(errmsg string) {
		logger.Println(errmsg)
		http.Error(w, "an internal server error occurred", http.StatusInternalServerError)
	}
	if !loggedIn(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	user, err := FetchUser(w, req)
	if err != nil {
		intErr("error fetching user")
		return
	}
	if !TutGameExist(user.Email) {
		tut := NewTutorial()
		tut.Add(user.Email)
	}

	// at this point everyone should have a game created or be at the login screen
	http.Redirect(w, req, "/tutorial", http.StatusSeeOther)
}

func tutorial(w http.ResponseWriter, req *http.Request) {
	if !loggedIn(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "tutorial.html", struct{}{})
}
