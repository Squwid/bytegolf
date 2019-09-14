package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Squwid/bytegolf/sess"
)

func login(w http.ResponseWriter, req *http.Request) {
	loggedIn, err := sess.LoggedIn(req)
	if err != nil {
		w.Write([]byte("error checking for login: " + err.Error()))
		return
	}
	if loggedIn {
		w.Write([]byte("you are already logged in"))
		return
	}
	var username = "frank"
	s, err := sess.Login(username)
	if err != nil {
		w.Write([]byte("error logging in:" + err.Error()))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "bg-session",
		Value: s.ID,
	})
	w.Write([]byte("you should now be logged in"))

	// _, err := req.Cookie("bg-session")
	// if err == nil {
	// 	// this does not guarentee they are logged in
	// 	w.Write([]byte("you are already logged in"))
	// 	return
	// }
	// uid := uuid.New().String()
	// // the username would need to come from a database on login from github
	// s := sess.Session{
	// 	ID:       uid,
	// 	Username: "something",
	// }
	// err = s.Put()
	// if err != nil {
	// 	w.Write([]byte("couldnt log in due to error:" + err.Error()))
	// 	return
	// }
	// http.SetCookie(w, &http.Cookie{
	// 	Name:  "bg-session",
	// 	Value: uid,
	// })
	// w.Write([]byte("you should now be logged in"))
}

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
	fmt.Println(time.Now().Local().Unix())
}
