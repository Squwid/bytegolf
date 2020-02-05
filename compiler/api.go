// this file is going to be all of the api calls to receive different types of executes for a specific user
// or for the overall leaderboards

package compiler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Squwid/bytegolf/firestore"
	"github.com/Squwid/bytegolf/sess"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Error reading all of code from compile post request: %v", err)
		// since this probably happened because some a-hole submitted super long code so send a
		// bad request back.
		// TODO: should i check the length here and see if it is too much memory for my little container
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var exe Execute
	err = json.Unmarshal(bs, &exe)
	if err != nil {
		log.Errorf("Error unmarshalling code from compile post request: %v", err)
		// this is probably a json formatting thing but im not even sure if that errors out and im not going
		// to send myself bad json so just return a bad error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: validate incoming execute requests

	resp, err := exe.Post(s)
	if err != nil {
		log.Errorf("Error compiling code from post request %v: %v", exe.HoleID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// once we get the response, write it to the api
	bs, err = json.Marshal(resp)
	if err != nil {
		log.Errorf("Error marshalling compile request for hole %v: %v", exe.HoleID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

// getQuestions returns a list of executes using a BGID. Currently gets all but should
// max out at some point
func getQuestions(BGID, holeID string, max int) ([]TotalStore, error) {
	ctx := context.Background()
	iter := firestore.Client.Collection("executes").Where("BGID", "==", BGID).Where("HoleID", "==", holeID).Limit(max).Documents(ctx)
	var exes = []TotalStore{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var exe TotalStore
		err = mapstructure.Decode(doc.Data(), &exe)
		if err != nil {
			return nil, err
		}
		exes = append(exes, exe)
	}
	return exes, nil
}
