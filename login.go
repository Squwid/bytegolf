package main

import (
	"net/http"

	"github.com/Squwid/bytegolf/sess"
)

func isLoggedIn(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	loggedIn, err := sess.LoggedIn(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !loggedIn {
		w.Write([]byte(`{"logged_in": false}`))
		return
	}
	w.Write([]byte(`{"logged_in": true}`))
}
