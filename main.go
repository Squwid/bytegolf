package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Squwid/bytegolf/bgaws"
	"github.com/Squwid/bytegolf/runner"
)

// CurrentGame is the current game of code golf
var CurrentGame Game
var tpl *template.Template
var logger *log.Logger
var currentSessions = map[string]session{}

// Errors
var (
	ErrGameFull = errors.New("game already has maximum amount of players")
)

// var currentGame = map[string]Game{} // maps a players name to a game

// Player struct that holds each players hole submissions
type Player struct {
	User          bgaws.User    // holds username, password, and role
	Scores        map[int]int64 // Scores holds each hole and what the player scored on it
	Correct       map[int]bool  // whether or not the player got the scores correct
	CorrectAmount int
}

// Game struct
type Game struct {
	ID             string
	Name           string
	Password       string
	CurrentPlayers int
	MaxPlayers     int
	Holes          int
	Difficulty     string
	StartedTime    time.Time
	Started        bool
	Players        []*Player
	Questions      map[int]bgaws.Question
}

// GolfResponse TODO: this structure needs to be removed at some point because we need anon structs eventually
type golfResponse struct {
	User     *bgaws.User
	Name     string
	LoggedIn bool
	Game     Game
	GameName string
	Question bgaws.Question
	Hole     int
}

type session struct {
	Username     string
	lastActivity time.Time
}

func init() {
	tpl = template.Must(template.ParseGlob("views/*"))
	logger = log.New(os.Stdout, "[bytegolf] ", log.Ldate|log.Ltime)
}

func main() {
	// host files
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	// handlers
	http.HandleFunc("/", index)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/currentgame/", current)
	http.HandleFunc("/logout", logout)

	// listen and serve
	http.Handle("/favicon.ico", http.NotFoundHandler())
	logger.Println("listening on port :6017")
	http.ListenAndServe(":6017", nil)
}

func createPlayer(user *bgaws.User) *Player {
	return &Player{
		User:          *user,
		Scores:        make(map[int]int64),
		Correct:       make(map[int]bool),
		CorrectAmount: 0,
	}
}

// AddGameUser adds a user to the current game
func (game *Game) AddGameUser(user *bgaws.User) error {
	if game.MaxPlayers == game.CurrentPlayers {
		return ErrGameFull
	}
	player := createPlayer(user)
	game.CurrentPlayers++ // add the player to list of players
	game.Players = append(game.Players, player)

	logger.Printf("%s added to game %s\n", user.Username, game.Name)
	logger.Printf("there are now %v people in game %s\n", game.CurrentPlayers, game.Name)
	return nil
}

// UserInGame checks to see if a user is in a specific game
func (game *Game) UserInGame(user *bgaws.User) bool {
	for _, p := range game.Players {
		if p.User.Username == user.Username {
			return true
		}
	}
	return false
}

// PlayerInGame checks to see if a certain player is in a game
func (game *Game) PlayerInGame(player *Player) bool {
	for _, p := range game.Players {
		if p.User.Username == player.User.Username {
			return true
		}
	}
	return false
}

// checks the response compared to the question TODO: Instead of a bool change this to something easier
func checkResponse(resp *runner.CodeResponse, q *bgaws.Question) bool {
	if strings.TrimSpace(strings.ToLower(resp.Output)) == strings.TrimSpace(strings.ToLower(q.Answer)) {
		return true
	}
	return false
}

func score(submission *runner.CodeSubmission, q *bgaws.Question) (p int64) {
	// TODO: now is just the length of the code, however i would like a better score system in the future
	return count(submission.Script)
}

func count(s string) int64 {
	var c int64
	for _, l := range s {
		if len(strings.TrimSpace(string(l))) == 0 {
			continue
		}
		c++
	}
	return c
}
