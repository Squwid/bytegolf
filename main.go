package main

import (
	"crypto/tls"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/Squwid/bytegolf/questions"
	_ "github.com/Squwid/bytegolf/users"
	"github.com/aws/aws-sdk-go/aws"
	awss "github.com/aws/aws-sdk-go/aws/session"
	"golang.org/x/crypto/acme/autocert"
)

var gitState = "abcdefg" // TODO: change this after testing if it works
var siteAddr = "https://bytegolf.io"
var tpl *template.Template
var awsSess *awss.Session

// Loggers
var (
	logger *log.Logger
)

type code struct {
	Show    bool
	Correct bool
	Bytes   int64
	Output  string
}

func init() {
	tpl = template.Must(template.ParseGlob("views/*"))
	logger = log.New(os.Stdout, "[bytegolf] ", log.Ldate|log.Ltime)
	awsSess = awss.Must(awss.NewSessionWithOptions(awss.Options{Config: aws.Config{Region: aws.String("us-east-1")}}))
	setGitClient()
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
	mux.HandleFunc("/play/", play)
	mux.HandleFunc("/submit/", submission)
	mux.HandleFunc("/holes/", holes)
	mux.HandleFunc("/login", gitLogin)
	mux.HandleFunc("/login/check", githubOAUTH)
	mux.HandleFunc("/account/", account)
	mux.HandleFunc("/leaderboards", leaderboards)

	/* ADMIN FUNCTIONS */
	mux.HandleFunc("/admin/", admin)
	mux.HandleFunc("/admin/archive/", archiveQuestion)
	mux.HandleFunc("/admin/deploy/", deployQuestion)
	mux.HandleFunc("/admin/refresh", refreshQuestions)
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
}
