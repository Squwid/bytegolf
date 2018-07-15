package main

import (
	"html/template"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("views/*"))
}

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	http.HandleFunc("/", index)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/dev", dev)
	http.HandleFunc("/currentgame", current)
	http.HandleFunc("/leaderboards", leaderboard)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":6017", nil)
}
