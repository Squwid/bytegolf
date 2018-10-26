package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Squwid/bytegolf/aws"
	"github.com/Squwid/bytegolf/runner"
)

/*
	TODO:
	DEV NOTES
	Currently need to implement the following before alpha
	* Deployable users.json that allows for users to be saved locally (just do user names on each login)
	* Remove comments from scoring
*/

// session id stored on players computer
// sessions stored to email
// (sessionID) -> (sessions) -> (email) -> (game) -> (players) -> player

// CurrentGame holds the current game, maybe support for more than one game in the future
var CurrentGame Game

var tpl *template.Template
var sessions = map[string]session{}

// Loggers
var (
	logger *log.Logger
	config *log.Logger
)

type session struct {
	Email        string
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
	go tutJanitor() // send off the janitor to always be running in a seperate thread
	// host files
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	// handlers
	http.HandleFunc("/", index)
	// http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/play", play)
	http.HandleFunc("/account", account)
	http.HandleFunc("/leaderboards", leaderboards)
	http.HandleFunc("/create", create)
	http.HandleFunc("/tutorial/create", tutCreator)
	http.HandleFunc("/tutorial", tutorial)

	// listen and serve
	http.Handle("/favicon.ico", http.NotFoundHandler())
	logger.Printf("listening on port :%s\n", Config.Port)
	http.ListenAndServe(":"+Config.Port, nil)
}

func checkResponse(resp *runner.CodeResponse, q *aws.Question) bool {
	if strings.TrimSpace(strings.ToLower(resp.Output)) == strings.TrimSpace(strings.ToLower(q.Answer)) {
		return true
	}
	return false
}

func score(sub *runner.CodeSubmission, q *aws.Question) int64 {
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
