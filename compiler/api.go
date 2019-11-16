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

const maxReturns = 20

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
		w.WriteHeader(http.StatusForbidden)
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
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS,GET")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodOptions {
		// this is for cors
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method == http.MethodGet {
		getExecutes(w, r, s)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var exe Execute
	err = json.Unmarshal(bs, &exe)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// change to secrets manager from environmental variables
	exe.ClientID = jdoodleClient.Client
	exe.ClientSecret = jdoodleClient.Secret

	resp, err := exe.Post(s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bs, err = json.Marshal(*resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(bs)
}

// getExecutes runs each time a hole is loaded to grab all of the responses, if none exist then a blank list
// will be returned meaning that the user has never submited a successful
func getExecutes(w http.ResponseWriter, r *http.Request, s *sess.Session) {
	hole := r.URL.Query().Get("hole")
	if hole == "" {
		// no hole on get means bad request
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	qs, err := getQuestions(s.BGID, hole, maxReturns)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if len(qs) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}
	bs, err := json.Marshal(qs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("error marshalling execute list: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	return
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
