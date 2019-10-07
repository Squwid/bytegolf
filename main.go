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
)

var siteAddr = "https://bytegolf.io"

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("dist/frontend/index.html"))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	http.Handle("/dist/frontend/", http.StripPrefix("/dist/frontend", http.FileServer(http.Dir("./dist/frontend"))))

	// handlers
	http.Handle("/", frontend("dist/frontend"))
	http.HandleFunc("/login/check", github.Oauth)
	http.HandleFunc("/login", github.Login)
	http.HandleFunc("/check", isLoggedIn)
	http.HandleFunc("/compile", compiler.Handler)
	http.ListenAndServe(":"+port, nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("byte golf api"))
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
