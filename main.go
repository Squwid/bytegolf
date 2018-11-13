package main

import (
	"crypto/tls"
	"fmt"
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

var siteAddr = "https://bytegolf.io"
var tpl *template.Template
var sessions = map[string]session{}
var questions = map[int]aws.Question{}

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
	questions = aws.GetQuestionsTemp(3)
	logger = log.New(os.Stdout, "[bytegolf] ", log.Ldate|log.Ltime)
	Config = SetupConfiguration(ParseConfiguration())
}

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	// handlers
	http.HandleFunc("/", index)
	// http.HandleFunc("/signup", signup)
	http.HandleFunc("/play/", play)
	http.HandleFunc("/holes/", holes)
	http.HandleFunc("/login", login)
	http.HandleFunc("/account", account)
	http.HandleFunc("/leaderboards", leaderboards)
	http.HandleFunc("/tutorial/create", tutCreator)
	http.HandleFunc("/tutorial", tutorial)
	http.HandleFunc("/.well-known/pki-validation", verify)

	// listen and serve
	http.Handle("/favicon.ico", http.NotFoundHandler())
	logger.Printf("listening on port :%s\n", "80")
	http.ListenAndServe(":"+"80", nil)
	// go http.ListenAndServe(":80", http.HandlerFunc(redirect))
	// srv := createServer()
	// srv.ListenAndServeTLS("", "") // TODO: need both certs
}

func verify(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("./4D4D9FD22C3F52BD286EAFD7DA81E122.txt")
	t.Execute(w, nil)
}

func redirect(w http.ResponseWriter, req *http.Request) {
	var target string
	if len(req.URL.RawQuery) > 0 {
		target = "?" + req.URL.RawQuery
	}
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

func createServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/play/", play)
	mux.HandleFunc("/holes/", holes)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/account", account)
	mux.HandleFunc("/leaderboards", leaderboards)
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS11,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         ":443",
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
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

// getHoleByLink retrieves an aws question from the questions and an error if one is not found
// with a matching link
func getHoleByLink(link string) (*aws.Question, error) {
	for _, hole := range questions {
		if hole.Link == link {
			return &hole, nil
		}
	}
	return nil, fmt.Errorf("could not find hole %s", link)
}
