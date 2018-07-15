package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Squwid/bytegolf/bgaws"
)

var tpl *template.Template
var currentSessions = map[string]session{} // sessionID : session
var currentGame *bgaws.Game
var logger *log.Logger

// GolfResponse TODO
type golfResponse struct {
	User     *bgaws.User
	Name     string
	LoggedIn bool
	Game     *bgaws.Game
	GameName string
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
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	http.HandleFunc("/", index)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/dev", dev)
	http.HandleFunc("/currentgame", current)
	http.HandleFunc("/leaderboards", leaderboard)

	http.Handle("/favicon.ico", http.NotFoundHandler())
	logger.Println("listening on port :6017")
	http.ListenAndServe(":6017", nil)
}
