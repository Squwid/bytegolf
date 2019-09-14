package main

import (
	"net/http"

	_ "github.com/Squwid/bytegolf/firestore"
)

var siteAddr = "https://bytegolf.io"

func init() {
	// rdsClient = redis.NewClient(&redis.Options{
	// 	Addr:     "redis:80",
	// 	Password: "",
	// 	DB:       0,
	// })
}

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	// handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("byte golf"))
	})
	http.HandleFunc("/login", login)
	http.HandleFunc("/check", isLoggedIn)
	// http.HandleFunc("/", index)
	// http.HandleFunc("/testing", sesss)
	// http.HandleFunc("/compile", compile)
	http.ListenAndServe(":80", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("byte golf api"))
}
