package main

import (
	"net/http"

	"github.com/Squwid/bytegolf/auth"
	_ "github.com/Squwid/bytegolf/firestore"
	"github.com/Squwid/bytegolf/globals"
	"github.com/Squwid/bytegolf/submissions"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var siteAddr = "https://bytegolf.io"

func main() {
	// getting the port here is essential when using google cloud run
	port := globals.Port()

	r := mux.NewRouter()

	// handlers
	r.HandleFunc("/login/check", auth.CallbackHandler)
	r.HandleFunc("/login", auth.LoginHandler)

	// r.HandleFunc("/api/holes", question.ListQuestionsHandler) // list all of the holes
	// r.HandleFunc("/api/hole", question.SingleHandler)         // list a single hole

	// Check if a user is logged in for frontend purposes
	r.HandleFunc("/api/profile/", auth.ProfileHandler).Methods("GET") // checks if a user is logged in

	// Compile request
	// r.HandleFunc("/api/compile", compiler.Handler)
	// r.HandleFunc("/api/submissions", compiler.SubmissionsHandler)

	/* SUBMISSION HANDLERS */
	r.HandleFunc("/api/submissions/best/{hole}", submissions.GetPlayersBestSubmission).Methods("GET")
	r.HandleFunc("/api/submissions/{hole}", submissions.GetLeaderboardForHole).Methods("GET")
	r.HandleFunc("/api/submission/{id}", submissions.GetSingleSubmission).Methods("GET")

	log.Infof("Starting container on port :%s", port)
	http.ListenAndServe(":"+port, r)
}
