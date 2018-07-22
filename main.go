package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Squwid/bytegolf/bgaws"
)

var (
	tpl             *template.Template
	currentGame     *Game
	logger          *log.Logger
	currentSessions = map[string]session{} // sessionID : session
	currentHoles    = map[string]string{}  // each player mapped to their current hole
)

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
	Players        []string
}

// GolfResponse TODO
type golfResponse struct {
	User     *bgaws.User
	Name     string
	LoggedIn bool
	Game     *Game
	GameName string
	Question *bgaws.Question
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
	http.HandleFunc("/logout", logout)

	http.Handle("/favicon.ico", http.NotFoundHandler())
	logger.Println("listening on port :6017")
	http.ListenAndServe(":6017", nil)
}

func addUserToCurrent(w http.ResponseWriter, user bgaws.User) error {
	gameCookie := &http.Cookie{
		Name:  "gameid",
		Value: currentGame.ID,
	}
	// set the user to the first hole
	holeCookie := &http.Cookie{
		Name:  "hole",
		Value: "1",
	}

	currentGame.CurrentPlayers++ // add the player to list of players
	currentGame.Players = append(currentGame.Players, user.Username)
	logger.Printf("%s added to game %s\n", user.Username, currentGame.Name)
	logger.Printf("there are now %v people in game %s\n", currentGame.CurrentPlayers, currentGame.Name)
	http.SetCookie(w, gameCookie)

	return nil
}
