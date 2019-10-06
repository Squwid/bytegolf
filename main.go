package main

import (
	"html/template"
	"net/http"
	"os"

	_ "github.com/Squwid/bytegolf/firestore"
	"github.com/Squwid/bytegolf/github"
)

var siteAddr = "https://bytegolf.io"

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("views/*"))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	// handlers
	http.HandleFunc("/", indexHTML)
	http.HandleFunc("/login/check", github.Oauth)
	http.HandleFunc("/login", github.Login)
	http.HandleFunc("/check", isLoggedIn)
	http.ListenAndServe(":"+port, nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("byte golf api"))
}
