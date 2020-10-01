package submissions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const maxSubs = 5

// GetSingleSubmission gets a single submission by using an id
func GetSingleSubmission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	l := log.WithFields(log.Fields{
		"Action": "GetSingleSubmission",
		"ID":     id,
	})
	l.Infof("Request to get single submission")

	// Get the submission with that ID
	sub, err := GetSubmissionByID(id)
	if err != nil {
		l.WithError(err).Errorf("Error getting submission by id")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if sub == nil {
		l.Warnf("Did not find submission")
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	l.WithField("Sub", sub).Infof("Got submission")

	// Only marshal short response
	bs, err := json.Marshal(sub.Short)
	if err != nil {
		l.WithError(err).Errorf("Error marshalling short submission")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(bs)
}

// GetLeaderboardForHole gets all of the short submissions for a single hole
func GetLeaderboardForHole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	hole := mux.Vars(r)["hole"]

	l := log.WithFields(log.Fields{
		"Action": "ListSubsForHole",
		"Hole":   hole,
	})
	l.Infof("Request to list submissions for single hole")

	// Make sure that the hole exists
	/*
		q, err := question.GetQuestion(hole)
		if err != nil {
			l.WithError(err).Errorf("Error getting hole")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if q == nil {
			l.Warnf("Hole was not found")
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
	*/

	l.Infof("Hole was found, now getting submissions")

	// Get all of the submissions for that hole
	subs, err := GetBestSubmissionsOnHole(hole, maxSubs)
	if err != nil {
		l.WithError(err).Errorf("Error getting best submissions for hole")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Got subs, transform to short submissions
	shortSubs := subs.ToShortSubmissions()

	// Marshal submissions
	bs, err := json.Marshal(shortSubs)
	if err != nil {
		l.WithError(err).Errorf("Error marshalling short submissions")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Count", fmt.Sprintf("%v", shortSubs))
	w.Write(bs)
}

// GetPlayersBestSubmission is the handler to return the player's best score on the hole
func GetPlayersBestSubmission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	hole := mux.Vars(r)["hole"]

	var bgid = "9abf282b-9d05-4b11-b4d4-3f8a8af9ea2f" // TODO: un-hardcode this because the player needs to be logged in

	l := log.WithFields(log.Fields{
		"Action": "GetPlayerBestSub",
		"Hole":   hole,
		"BGID":   bgid,
	})
	l.Infof("Request to get player's best sub")

	// Get the players best submission
	sub, err := GetBestPlayerSubmission(bgid, hole)
	if err != nil {
		l.WithError(err).Errorf("Error getting player's best submission")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if sub == nil {
		l.Warnf("No submissions found for this hole")
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	l.WithField("Sub", sub).Infof("Got players best submission!")

	// Marshal the submission and return it
	bs, err := json.Marshal(sub)
	if err != nil {
		l.WithError(err).Errorf("Error marshalling submission")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(bs)
}
