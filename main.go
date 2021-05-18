package main

import (
	"net/http"

	"github.com/Squwid/bytegolf/auth"
	_ "github.com/Squwid/bytegolf/firestore"
	"github.com/Squwid/bytegolf/globals"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var siteAddr = "https://bytegolf.io"

func main() {
	// getting the port here is essential when using google cloud run
	port := globals.Port()

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })

	r.HandleFunc("/login/check", auth.CallbackHandler)
	r.HandleFunc("/login", auth.LoginHandler)

	// Check if a user is logged in for frontend purposes

	r.HandleFunc("/api/profile/{bgid}", auth.ShowProfile).Methods("GET") // checks if a user is logged in
	r.HandleFunc("/api/claims", auth.ShowClaims).Methods("GET")          // Returns a user's claims

	log.Infof("Starting container on port :%s", port)
	http.ListenAndServe(":"+port, r)
}
