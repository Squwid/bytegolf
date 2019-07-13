package main

import (
	"crypto/tls"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Squwid/bytegolf/database"
	_ "github.com/Squwid/bytegolf/questions"
	"github.com/Squwid/bytegolf/users"
	"github.com/aws/aws-sdk-go/aws"
	awss "github.com/aws/aws-sdk-go/aws/session"
	"golang.org/x/crypto/acme/autocert"
)

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
}

func main() {
	if database.InProd() {
		mux := http.NewServeMux()
		certManager := autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache("cert-cache"),
		}

		mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

		// handlers
		mux.HandleFunc("/", index)
		mux.HandleFunc("/play/", play)
		mux.HandleFunc("/submit/", submission)
		mux.HandleFunc("/holes/", holes)
		mux.HandleFunc("/login", users.GitLogin)
		mux.HandleFunc("/login/check", users.GithubOAUTH)
		mux.HandleFunc("/account/", account)
		mux.HandleFunc("/leaderboards", leaderboards)

		/* ADMIN FUNCTIONS */
		mux.HandleFunc("/admin/", admin)
		mux.HandleFunc("/admin/archive/", archiveQuestion)
		mux.HandleFunc("/admin/deploy/", deployQuestion)
		mux.HandleFunc("/admin/addhole", createQuestion)
		mux.HandleFunc("/admin/delete/", deletehole)
		mux.HandleFunc("/admin/holes", adminholes)

		// listen and serve
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
	} else {
		http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

		// handlers
		http.HandleFunc("/", index)
		http.HandleFunc("/play/", play)
		http.HandleFunc("/submit/", submission)
		http.HandleFunc("/holes/", holes)
		http.HandleFunc("/login", users.GitLogin)
		http.HandleFunc("/login/check", users.GithubOAUTH)
		http.HandleFunc("/account/", account)
		http.HandleFunc("/leaderboards", leaderboards)

		/* ADMIN FUNCTIONS */
		http.HandleFunc("/admin/", admin)
		http.HandleFunc("/admin/archive/", archiveQuestion)
		http.HandleFunc("/admin/deploy/", deployQuestion)
		http.HandleFunc("/admin/addhole", createQuestion)
		http.HandleFunc("/admin/delete/", deletehole)
		http.HandleFunc("/admin/holes", adminholes)
		http.ListenAndServe(":80", nil)
	}
}
