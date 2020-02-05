package main

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/Squwid/bytegolf/compiler"
	_ "github.com/Squwid/bytegolf/firestore"
	"github.com/Squwid/bytegolf/github"
	question "github.com/Squwid/bytegolf/questions"
)

var siteAddr = "https://bytegolf.io"

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("dist/frontend/index.html"))
}

func main() {
	// getting the port here is essential when using google cloud run
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	// host file server for frontend assets
	http.Handle("/dist/frontend/", http.StripPrefix("/dist/frontend", http.FileServer(http.Dir("./dist/frontend"))))

	// handlers
	http.Handle("/", frontend("dist/frontend"))
	http.HandleFunc("/login/check", github.Oauth)
	http.HandleFunc("/login", github.Login)
	http.HandleFunc("/api/holes", question.ListQuestionsHandler) // list all of the holes
	http.HandleFunc("/api/hole", question.SingleHandler)         // list a single hole
	http.HandleFunc("/api/user", isLoggedIn)                     // checks if a user is logged in
	http.HandleFunc("/compile", compiler.Handler)
	http.HandleFunc("/api/submissions", compiler.SubmissionsHandler)
	http.ListenAndServe(":8080", nil)
}

func frontend(dir string) http.Handler {
	handler := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		p := req.URL.Path
		if strings.Contains(p, ".") || p == "/" {
			handler.ServeHTTP(w, req)
			return
		}
		http.ServeFile(w, req, path.Join(dir, "/index.html"))
	})
}
