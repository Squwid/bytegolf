package main

import (
	"crypto/tls"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/Squwid/bytegolf/questions"
	"github.com/aws/aws-sdk-go/aws"
	awss "github.com/aws/aws-sdk-go/aws/session"
)

/*
	DEV NOTES
	Currently need to implement the following before alpha
*/

// session id stored on players computer
// sessions stored to email

// Maximum amount of holes possible at a single time
const maxHoles = 9

var siteAddr = "https://bytegolf.io"
var tpl *template.Template
var sessions = map[string]session{}
var awsSess *awss.Session

// Loggers
var (
	logger *log.Logger
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
	// qs = questions.ToMap(questions.GetLocalQuestions())
	logger = log.New(os.Stdout, "[bytegolf] ", log.Ldate|log.Ltime)
	awsSess = awss.Must(awss.NewSessionWithOptions(awss.Options{Config: aws.Config{Region: aws.String("us-east-1")}}))
}

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	// handlers
	http.HandleFunc("/", index)
	// http.HandleFunc("/signup", signup)
	http.HandleFunc("/play/", play)
	http.HandleFunc("/submit/", submission)
	http.HandleFunc("/holes/", holes)
	http.HandleFunc("/login", login)
	http.HandleFunc("/account", account)
	http.HandleFunc("/leaderboards", leaderboards)

	/* ADMIN FUNCTIONS */
	http.HandleFunc("/admin/", admin)
	http.HandleFunc("/admin/archive/", archiveQuestion)
	http.HandleFunc("/admin/deploy/", deployQuestion)
	http.HandleFunc("/admin/refresh", refreshQuestions)

	// http.HandleFunc("/admin/delete/", deletehole)
	http.HandleFunc("/admin/addhole", createQuestion)
	http.HandleFunc("/admin/adduser", createUser)
	http.HandleFunc("/admin/holes", adminholes)

	// listen and serve
	http.Handle("/favicon.ico", http.NotFoundHandler())
	logger.Printf("listening on port :%s\n", "80")
	http.ListenAndServe(":80", nil)
	// go http.ListenAndServe(":80", http.HandlerFunc(redirect))
	// srv := createServer()
	// srv.ListenAndServeTLS("", "") // TODO: need both certs
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
