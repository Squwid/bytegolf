package compiler

import (
	"context"
	"encoding/json"
	"net/http"

	"cloud.google.com/go/firestore"
	fs "github.com/Squwid/bytegolf/firestore"
	"github.com/Squwid/bytegolf/sess"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

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
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// if the query string best is true, return the users best score
	if r.URL.Query().Get("best") == "true" {
		bestLong, err := GetBestSubmission(s.BGID, hole)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if bestLong == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		best := bestLong.TransformToShort()

		bs, err := json.Marshal(best)
		if err != nil {
			log.Errorf("Error marshalling best sub %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(bs)
		return
	}

	// if subid is not none, then get a short submi
	subID := r.URL.Query().Get("id")
	if subID != "" {
		sub, err := GetSingleSubmission(subID)
		if err != nil {
			log.Errorf("Error getting single submission: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if sub == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// the sub is good so transform it to return a short submission + the
		ss := sub.TransformToShort()
		bs, err := json.Marshal(struct {
			ShortSubmission
			Script string `json:"script"`
		}{
			ShortSubmission: ss,
			Script:          sub.Exe.Script,
		})
		if err != nil {
			log.Errorf("Error getting a submission for %s with id %s", s.BGID, subID)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(bs)
		return
	}

	// get request so get all of the requests for that user
	subs, err := ListShortSubmissions(s.BGID, hole, maxReturns)
	if err != nil {
		log.Errorf("Error listing short submissions: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(subs) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))
		return
	}
	bs, err := json.Marshal(subs)
	if err != nil {
		log.Errorf("Error marshalling submission list: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

// GetSingleSubmission gets a single submission by using an id
func GetSingleSubmission(id string) (*FullSubmission, error) {
	ctx := context.Background()
	iter := fs.Client.Collection("executes").Where("UUID", "==", id).Limit(1).Documents(ctx)

	var sub FullSubmission

	doc, err := iter.Next()
	if err == iterator.Done {
		// no result was found for that execution
		log.Warnf("Request to get single submission %s was not found", id)
		return nil, nil
	}
	if err != nil {
		log.Errorf("Error getting single submission %s: %v", id, err)
		return nil, err
	}

	err = mapstructure.Decode(doc.Data(), &sub)
	if err != nil {
		log.Errorf("Error decoding single submission %s: %v", id, err)
		return nil, err
	}

	return &sub, nil
}

// ListShortSubmissions lists a users short submissions
func ListShortSubmissions(bgid, hole string, max int) ([]ShortSubmission, error) {
	ctx := context.Background()
	iter := fs.Client.Collection("executes").Where("BGID", "==", bgid).
		Where("HoleID", "==", hole).Limit(max).Documents(ctx)

	var ss = []ShortSubmission{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Errorf("Error listing short submissions %s %v: %v", bgid, hole, err)
			return nil, err
		}

		var sub FullSubmission
		err = mapstructure.Decode(doc.Data(), &sub)
		if err != nil {
			log.Errorf("Error decoding short submission: %s %v: %v", bgid, hole, err)
			return nil, err
		}

		short := sub.TransformToShort()
		ss = append(ss, short)
	}
	return ss, nil
}

// GetBestSubmission gets the best submission for a specific user
func GetBestSubmission(bgid, hole string) (*FullSubmission, error) {
	ctx := context.Background()
	iter := fs.Client.Collection("executes").Where("BGID", "==", bgid).
		Where("HoleID", "==", hole).Where("Correct", "==", true).Limit(1).
		OrderBy("Length", firestore.Asc).Documents(ctx)

	var sub FullSubmission
	doc, err := iter.Next()
	if err == iterator.Done {
		log.Warnf("Request to get best submission %s (%s) was not found", bgid, hole)
		return nil, nil
	}
	if err != nil {
		log.Errorf("Error getting best submission %s (%s): %v", bgid, hole, err)
		return nil, err
	}

	err = mapstructure.Decode(doc.Data(), &sub)
	if err != nil {
		log.Errorf("Error decoding best submission %s (%s): %v", bgid, hole, err)
		return nil, err
	}

	return &sub, nil
}

// GetBestSubmissionsOnHole gets the best submissions on a hole, this function is used for leaderboards
func GetBestSubmissionsOnHole(hole string, max int) ([]FullSubmission, error) {
	log.Infof("Request to list %v leaders on hole %s", max, hole)
	ctx := context.Background()
	iter := fs.Client.Collection("executes").Where("HoleID", "==", hole).Where("Correct", "==", true).
		OrderBy("Length", firestore.Asc).Documents(ctx)

	// contains is a function to see if the bgid is already in the slice
	contains := func(ss []FullSubmission, bgid string) bool {
		for _, s := range ss {
			if s.BGID == bgid {
				return true
			}
		}
		return false
	}

	var subs = []FullSubmission{}
	for {
		// since this is dynamic and it keeps reading,
		if len(subs) >= max {
			break
		}

		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Errorf("Error listing leaders for hole (%s): %v", hole, err)
			return nil, err
		}

		var sub FullSubmission
		err = mapstructure.Decode(doc.Data(), &sub)
		if err != nil {
			log.Errorf("Error decoding leaderboard submission for hole (%s): %v", hole, err)
			return nil, err
		}

		// only append to the slice if the user's bgid is not already there
		if !contains(subs, sub.BGID) {
			subs = append(subs, sub)
		}
	}
	log.Infof("Request to list %v best players on hole %s was successful and returned %v leaders", max, hole, len(subs))
	return subs, nil
}
