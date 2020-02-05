package compiler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	fs "github.com/Squwid/bytegolf/firestore"
	"github.com/Squwid/bytegolf/sess"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

/*
	type submission struct {
		ID            string    `json:"id"`
		Correct       bool      `json:"correct"`
		Language      string    `json:"language"`
		Score         int       `json:"score"`
		Script        string    `json:"script"`
		SubmittedTime time.Time `json:"submitted_time"`
	}
*/

// ShortSubmission is a submission with no script
type ShortSubmission struct {
	ID            string    `json:"id"`
	Correct       bool      `json:"correct"`
	Language      string    `json:"language"`
	Score         int       `json:"score"`
	SubmittedTime time.Time `json:"submitted_time"`
}

// Submission is the type that is returned by the submission
type Submission struct {
	ShortSubmission
	Script string `json:"script"`
}

// SubmissionsHandler handles all of the submissions api stuff
func SubmissionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodOptions {
		// Cors stuff
		w.WriteHeader(http.StatusOK)
		return
	}

	// check hole first because all other logic would be obselete
	hole := r.URL.Query().Get("hole")
	if hole == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// if the user is not logged in return a 503, they have to be signed in to see their response
	loggedIn, s, err := sess.LoggedIn(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("error checking if a user is signed in for submissions: %v", err)
		return
	}
	if !loggedIn {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "unauthorized"}`))
		return
	}
	if s == nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("error: session was blank")
		return
	}

	// Only accept get methods
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// if the query string best is true, return the users best score
	if r.URL.Query().Get("best") == "true" {
		getBestSubmission(w, s, hole)
		return
	}

	subID := r.URL.Query().Get("id")
	if subID != "" {
		getSingleSubmission(w, s, subID, hole)
		return
	}

	// get request so get all of the requests for that user
	listSubmissions(w, s, hole)
}

func getSingleSubmission(w http.ResponseWriter, s *sess.Session, subID, hole string) {
	ctx := context.Background()
	iter := fs.Client.Collection("executes").Where("BGID", "==", s.BGID).Where("HoleID", "==", hole).
		Where("UUID", "==", subID).Limit(1).Documents(ctx)

	var sub *Submission
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Errorf("Error getting single submission %s for %s on %s: %v", subID, s.BGID, hole, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var ts TotalStore
		err = mapstructure.Decode(doc.Data(), &ts)
		if err != nil {
			log.Errorf("Error decoding single submission %s for %s on %s: %v", subID, s.BGID, hole, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		sub = &Submission{
			Script: ts.Exe.Script,
			ShortSubmission: ShortSubmission{
				ID:            ts.UUID,
				Correct:       ts.Correct,
				Language:      ts.Exe.Language,
				Score:         ts.Length,
				SubmittedTime: ts.SubmittedTime,
			},
		}
	}

	if sub == nil {
		w.WriteHeader(http.StatusNotFound)
		log.Warnf("Request to single submission but it wasnt found %s for %s on %s", subID, s.BGID, hole)
		return
	}

	bs, err := json.Marshal(sub)
	if err != nil {
		log.Errorf("Error marshalling non nil best sub for %s on %s: %v", s.BGID, hole, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// WE DID IT
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func getBestSubmission(w http.ResponseWriter, s *sess.Session, hole string) {
	ctx := context.Background()
	iter := fs.Client.Collection("executes").Where("BGID", "==", s.BGID).Where("HoleID", "==", hole).
		Where("Correct", "==", true).
		Limit(1).OrderBy("Length", firestore.Asc).Documents(ctx)

	var sub *Submission
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Errorf("Error getting best submission for %s on %s: %v", s.BGID, hole, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var ts TotalStore
		err = mapstructure.Decode(doc.Data(), &ts)
		if err != nil {
			log.Errorf("Error decoding best submission for %s on %s: %v", s.BGID, hole, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		sub = &Submission{
			Script: ts.Exe.Script,
			ShortSubmission: ShortSubmission{
				ID:            ts.UUID,
				Correct:       ts.Correct,
				Language:      ts.Exe.Language,
				Score:         ts.Length,
				SubmittedTime: ts.SubmittedTime,
			},
		}
	}

	// check if the submission is nil, if so return a 404 so i can handle it on the frontend
	if sub == nil {
		w.WriteHeader(http.StatusNotFound)
		log.Warnf("Request to find best score for %s on %s but not found", s.BGID, hole)
		return
	}

	// the sub is NOT nil so marshal it and return it to the user
	bs, err := json.Marshal(sub)
	if err != nil {
		log.Errorf("Error marshalling non nil best sub for %s on %s: %v", s.BGID, hole, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// WE DID IT
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

// listSubmissions runs each time a hole is loaded to grab all of the responses, if none exist then a blank list
// will be returned meaning that the user has never submited a successful
func listSubmissions(w http.ResponseWriter, s *sess.Session, hole string) {
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

	// Getting all past submissions will then be parsed to a new structure to make sure that
	// the user doesnt see anything that they shouldnt
	var submissions = []ShortSubmission{}

	for _, q := range qs {
		// no script here because no need to pass that much data back and forth if
		// its not going to be used that often
		submissions = append(submissions, ShortSubmission{
			ID:            q.UUID,
			Correct:       q.Correct,
			Language:      q.Exe.Language,
			Score:         len(q.Exe.Script),
			SubmittedTime: q.SubmittedTime,
		})
	}

	bs, err := json.Marshal(submissions)
	if err != nil {
		log.Errorf("Error unmarshalling new short submissions: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Infof("Got %v responses for %v hole %v", len(submissions), s.BGID, hole)
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	return
}
