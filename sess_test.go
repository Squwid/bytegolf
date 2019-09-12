package main

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestSessions(t *testing.T) {
	sess := generateSession()
	err := sess.Add(time.Minute * 20)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	time.Sleep(1 * time.Second)
	in, s := authenticated(sess.ID)
	if !in {
		log.Errorf("Expected in to be true but got false\n")
		return
	}
	if sess.ID != s.ID {
		log.Errorf("Expected session id to be the same but got %s", s.ID)
	}
}
