package main

import (
	"net/http"

	"github.com/Squwid/bytegolf/sess"
)

func isLoggedIn(w http.ResponseWriter, req *http.Request) {
	loggedIn, err := sess.LoggedIn(req)
	if err != nil {
		w.Write([]byte("error trying to check login " + err.Error()))
		return
	}
	if !loggedIn {
		w.Write([]byte("not logged in"))
		return
	}
	w.Write([]byte("logged in"))
	// 1568570210 timeout: 1568483866 actual
	// fmt.Println(time.Now().Local().Unix())
}
