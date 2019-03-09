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
	"golang.org/x/crypto/acme/autocert"
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
	mux := http.NewServeMux()
	certManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache("certs"),
	}

	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	// handlers
	mux.HandleFunc("/", index)
	// http.HandleFunc("/signup", signup)
	mux.HandleFunc("/play/", play)
	mux.HandleFunc("/submit/", submission)
	mux.HandleFunc("/holes/", holes)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/account", account)
	mux.HandleFunc("/leaderboards", leaderboards)

	/* ADMIN FUNCTIONS */
	mux.HandleFunc("/admin/", admin)
	mux.HandleFunc("/admin/archive/", archiveQuestion)
	mux.HandleFunc("/admin/deploy/", deployQuestion)
	mux.HandleFunc("/admin/refresh", refreshQuestions)

	// mux.HandleFunc("/admin/delete/", deletehole)
	mux.HandleFunc("/admin/addhole", createQuestion)
	mux.HandleFunc("/admin/adduser", createUser)
	mux.HandleFunc("/admin/holes", adminholes)

	// listen and serve
	mux.Handle("/favicon.ico", http.NotFoundHandler())
	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         ":443",
		Handler:      mux,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}
	go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
	server.ListenAndServeTLS("", "")

	// logger.Printf("listening on port :%s\n", "80")
	// http.ListenAndServe(":80", nil)
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
