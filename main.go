package main

import (
	"net/http"

	"github.com/Squwid/bytegolf/auth"
	"github.com/Squwid/bytegolf/compiler"
	"github.com/Squwid/bytegolf/compiler/jdoodle"
	_ "github.com/Squwid/bytegolf/db"
	"github.com/Squwid/bytegolf/globals"
	"github.com/Squwid/bytegolf/holes"
	"github.com/Squwid/bytegolf/profiles"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	// getting the port here is essential when using google cloud runa
	port := globals.Port()

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })

	r.HandleFunc("/login/check", auth.CallbackHandler)
	r.HandleFunc("/login", auth.LoginHandler)

	r.HandleFunc("/api/holes", holes.ListHoles).Methods("GET")
	r.HandleFunc("/api/hole/{id}", holes.GetHole).Methods("GET")

	r.HandleFunc("/api/submit/{hole}", jdoodle.SubmissionHandler).Methods("POST")
	r.HandleFunc("/api/leaderboards", compiler.LeaderboardHandler).Methods("GET")

	r.HandleFunc("/api/profile/{id}", profiles.GetProfile).Methods("GET") // checks if a user is logged in
	r.HandleFunc("/api/claims", auth.ShowClaims).Methods("GET")           // Returns a user's claims and see if they are logged in

	log.Infof("Starting container on port :%s", port)
	http.ListenAndServe(":"+port, loggedIn(cors(r)))
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			w.WriteHeader(200)
			return
		}
		next.ServeHTTP(w, r)
	})
}
