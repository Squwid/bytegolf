package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var rdsClient *redis.Client

type session struct {
	// placeholder for now
	ID       string `json:"id"`
	Username string `json:"username"`
}

// generateSession creates a new session
func generateSession() (sess session) {
	sess.ID = uuid.New().String()
	return
}

func (s session) Add(dur time.Duration) error {
	return rdsClient.Set(s.ID+"-sess", s, time.Hour*48).Err()
}

// authenticated checks to see if a user is logged in
// if they are logged in it returns a session as well, ONLY if the bool is true
func authenticated(id string) (bool, *session) {
	val, err := rdsClient.Get(id + "-sess").Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		log.Errorf("error authenticating request %s", id)
		return false, nil
	}
	var sess session
	err = json.Unmarshal([]byte(val), &sess)
	if err != nil {
		log.Errorf("error unmarshalling request: %v", err)
		return false, nil
	}
	return true, &sess
}

func remove(id string) {
	log.Infoln("remove:", rdsClient.Del(id+"-sess").Err())
}

func sess(w http.ResponseWriter, req *http.Request) {
	if req.URL.Query().Get("login") != "" {
		sess := generateSession()
		err := sess.Add(time.Minute * 5)
		if err != nil {
			w.Write([]byte("error logging in"))
			return
		}
	}
	// time.Sleep(1 * time.Minute)

	cookie, err := req.Cookie("bgsess")
	if err != nil {
		w.Write([]byte("you are not logged in"))
		return
	}

	if req.URL.Query().Get("logout") != "" {
		remove(cookie.Value)
		w.Write([]byte("just tried to log out"))
		return
	}
	log.Printf("Logged in. Cookie: %v", cookie.Value)
	w.Write([]byte("logged in, cookie: " + cookie.Value))
}
