package main

import (
	"fmt"
	"net/http"

	"github.com/Squwid/bytegolf/github"
	"github.com/Squwid/bytegolf/sess"
)

func isLoggedIn(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	loggedIn, s, err := sess.LoggedIn(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !loggedIn {
		w.Write([]byte(`{"logged_in": false}`))
		return
	}

	user, err := github.RetreiveUser(s.BGID)
	if err != nil {
		w.Write([]byte("error: " + err.Error()))
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"logged_in": true, "username": "%s"}`, user.Username)))
}
