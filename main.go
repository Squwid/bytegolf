package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var tpl *template.Template

func main() {
	r := newRouter()
	http.ListenAndServe(":6017", r)
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	staticFileDir := http.Dir("./assets")
	staticFileHander := http.StripPrefix("/assets/", http.FileServer(staticFileDir))
	r.PathPrefix("/assets").Handler(staticFileHander).Methods("GET")

	// Handler Function Declaration
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/golf", contentHandler).Methods("GET")
	r.HandleFunc("/leaderboards", lbHandler).Methods("GET")
	return r
}

// Index handler
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseGlob("views/*")
	if err != nil {
		log.Fatalf("error exeucting template %v", err)
	}

	if err = tpl.ExecuteTemplate(w, "index.html", nil); err != nil {
		log.Fatalln("error executing template %v", err)
	}
}

// Leaderboard Handler
func lbHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseGlob("views/*")
	if err != nil {
		log.Fatalf("error exeucting template %v", err)
	}

	if err = tpl.ExecuteTemplate(w, "leaderboards.html", nil); err != nil {
		log.Fatalln("error executing template %v", err)
	}
}

// Content (new game and current games)
func contentHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseGlob("views/*")
	if err != nil {
		log.Fatalf("error exeucting template %v", err)
	}

	if err = tpl.ExecuteTemplate(w, "golf.html", nil); err != nil {
		log.Fatalln("error executing template %v", err)
	}
}
