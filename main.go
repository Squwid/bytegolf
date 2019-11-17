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
	http.HandleFunc("/holes", question.Handler)
	http.HandleFunc("/hole", question.SingleHandler)
	http.HandleFunc("/login", github.Login)
	http.HandleFunc("/check", isLoggedIn)
	http.HandleFunc("/compile", compiler.Handler)
	http.ListenAndServe(":"+port, nil)
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
