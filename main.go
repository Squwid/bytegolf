package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Squwid/bytegolf/aws"
	"github.com/Squwid/bytegolf/questions"
	"github.com/Squwid/bytegolf/runner"
)

/*
	TODO:
	DEV NOTES
	Currently need to implement the following before alpha
	* Spacing on questions
	* Remove comments from scoring
*/

// CurrentGame is the current game of code golf
var CurrentGame Game
var tpl *template.Template
var logger *log.Logger
var currentSessions = map[string]session{}
var games = map[string]*Game{}
var users []*aws.User

// var currentGame = map[string]Game{} // maps a players name to a game

// Player struct that holds each players hole submissions
type Player struct {
	User         aws.User      // holds username, password, and role
	Scores       map[int]int64 // Scores holds each hole and what the player scored on it
	Correct      map[int]bool  // whether or not the player got the scores correct
	Output       map[int]string
	HolesCorrect int
	TotalScore   int64
	Average      float64
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
	Questions      map[int]questions.Question
	Owner          *Player
	Leaderboard    struct {
		Winning      *Player
		OtherPlayers []*Player
	}
	Over bool
}

type session struct {
	Username     string
	lastActivity time.Time
}

type code struct {
	Show    bool
	Correct bool
	Bytes   int64
	Output  string
}

func init() {
	tpl = template.Must(template.ParseGlob("views/*"))
	logger = log.New(os.Stdout, "[bytegolf] ", log.Ldate|log.Ltime)
	Config = SetupConfiguration(ParseConfiguration())
}

func main() {
	// host files
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	// handlers
	http.HandleFunc("/", index)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/currentgame/", current)
	http.HandleFunc("/master", master)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/profile", profile)
	http.HandleFunc("/rules", rules)
	http.HandleFunc("/leaderboard", leaderboard)
	http.HandleFunc("/admin", admin)

	// listen and serve
	http.Handle("/favicon.ico", http.NotFoundHandler())
	logger.Printf("listening on port :%s\n", Config.Port)
	http.ListenAndServe(":"+Config.Port, nil)
}

func createPlayer(user *aws.User) *Player {
	return &Player{
		User:         *user,
		Scores:       make(map[int]int64),
		Correct:      make(map[int]bool),
		Output:       make(map[int]string),
		HolesCorrect: 0,
		TotalScore:   0,
		Average:      0.0,
	}
}

func checkResponse(resp *runner.CodeResponse, q *questions.Question) bool {
	if strings.TrimSpace(strings.ToLower(resp.Output)) == strings.TrimSpace(strings.ToLower(q.Answer)) {
		return true
	}
	return false
}

func checkCorrect(hole int, p *Player) code {
	var c code
	c.Show = true
	if p.Correct[hole] {
		c.Correct = true
		c.Bytes = p.Scores[hole]
		c.Output = p.Output[hole]
	} else {
		c.Correct = false
		c.Output = p.Output[hole]
	}
	return c
}

func score(sub *runner.CodeSubmission, q *questions.Question) int64 {
	// TODO: now is just the length of the code, however i would like a better score system in the future
	return count(sub.Script)
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
