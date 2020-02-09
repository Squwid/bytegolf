// this file is going to be all of the api calls to receive different types of executes for a specific user
// or for the overall leaderboards

package compiler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Squwid/bytegolf/question"
	"github.com/Squwid/bytegolf/sess"
	log "github.com/sirupsen/logrus"
)

const maxReturns = 100

// Handler is the rest api function handler for golang
func Handler(w http.ResponseWriter, r *http.Request) {
	// you can try this password but it wont work
	loggedIn, s, err := sess.LoggedIn(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("error checking to see if a user is signed in: %v", err)
		return
	}
	if !loggedIn {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf(`{"error": "unauthorized"}`)))
		return
	}
	if s == nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("error: session was blank")
		return
	}
	// the user is logged in

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodOptions {
		// this is for cors
		w.WriteHeader(http.StatusOK)
		return
	}

	// only accept post methods
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Error reading all of code from compile post request: %v", err)
		// since this probably happened because some a-hole submitted super long code so send a bad request back.
		// TODO: should i check the length here and see if it is too much memory for my little container
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// get the users execute code and make sure the hole exists, then send the test to the compiler
	var exe Execute
	err = json.Unmarshal(bs, &exe)
	if err != nil {
		log.Errorf("Error unmarshalling code from compile post request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "request body too long"}`))
		return
	}

	// check that the hole exists
	q, err := question.GetQuestion(exe.HoleID)
	if err != nil {
		log.Errorf("Error getting question from execute on id %v: %v", exe.HoleID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if q == nil {
		log.Errorf("Execution for question %s does not exist", exe.HoleID)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "hole not found"}`))
		return
	}

	// send the question to the compiler if the question exists and is not null
	fullSub, err := exe.RunTests(*q, s.BGID)
	if err != nil {
		log.Errorf("Error compiling hole %s for player %s", q.ID, s.BGID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	shortSub := fullSub.TransformToShort()

	// once we get the response, write it to the api
	bs, err = json.Marshal(shortSub)
	if err != nil {
		log.Errorf("Error marshalling compile request for hole %v: %v", exe.HoleID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}
