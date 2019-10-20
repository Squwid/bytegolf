// this file is going to be all of the api calls to receive different types of executes for a specific user
// or for the overall leaderboards

package compiler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Squwid/bytegolf/firestore"
	"github.com/Squwid/bytegolf/sess"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

const maxReturns = 20

// GetUsersRecentScores gets the users recent scores by using their cookie and getting the last
// amount of holes. It gets the hole from the query param 'hole'

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
