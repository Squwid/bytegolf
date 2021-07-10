package compiler

import (
	"encoding/json"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/Squwid/bytegolf/auth"
	"github.com/Squwid/bytegolf/db"
	"github.com/Squwid/bytegolf/models"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const subLimit = 20

// ListSubmissions needs to be able to list submissions for ONLY the logged in user
// Possible query strings:
//     "hole": Optional query string just to get submissions for a single hole
func ListSubmissions(w http.ResponseWriter, r *http.Request) {
	log := logrus.WithField("Action", "ListSubmissions")

	claims := auth.LoggedIn(r)
	if claims == nil {
		log.Infof("User not authenticated")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log = log.WithField("User", claims.BGID)

	query := db.SubmissionsCollection().Where("BGID", "==", claims.BGID).OrderBy("SubmittedTime", firestore.Desc).Limit(subLimit)

	if hole := r.URL.Query().Get("hole"); hole != "" {
		query = query.Where("HoleID", "==", hole)
		log = log.WithField("Hole", hole)
	}

	subs, err := SubmissionsQuery(query)
	if err != nil {
		log.WithError(err).Errorf("Error getting submissions")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Map submissions to short submissions for frontend
	// TODO: Convert this to use goroutines to make submissions concurrent
	var shortSubs = make([]ShortSubmission, len(subs))
	for i := 0; i < len(subs); i++ {
		ss, err := subs[i].ShortSub()
		if err != nil {
			log.WithError(err).Errorf("Error converting short sub")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		shortSubs[i] = *ss
	}

	log.WithField("Submissions", len(subs)).Infof("Got submissions")

	bs, _ := json.Marshal(shortSubs)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}

// GetSubmission gets a single submission if the users BGID is the same as the submission requested
func GetSubmission(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log := logrus.WithFields(logrus.Fields{
		"Action": "GetSubmission",
		"ID":     id,
	})

	claims := auth.LoggedIn(r)
	if claims == nil {
		log.Infof("User not authenticated")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log = log.WithField("User", claims.BGID)

	getter := models.NewGet(db.SubmissionsCollection().Doc(id), nil)
	doc, err := db.Get(getter)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			w.WriteHeader(http.StatusNotFound)
			log.Warnf("Sub not found")
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		log.WithError(err).Errorf("Error getting sub")
		return
	}

	var sub SubmissionDB
	if err := mapstructure.Decode(doc, &sub); err != nil {
		log.WithError(err).Errorf("Error decoding submission")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Make sure BGID matches submission
	if sub.BGID != claims.BGID {
		log.Infof("%s on sub doesnt match %s", sub.BGID, claims.BGID)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fs, err := sub.FullSub()
	if err != nil {
		log.WithError(err).Errorf("Error getting full submission")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Infof("Got full submission")

	bs, _ := json.Marshal(fs)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}

// SubmissionsQuery is a query wrapper to add results to a cache for when a GET is called at a later time
// TODO: add a cache
func SubmissionsQuery(query firestore.Query) ([]SubmissionDB, error) {
	docs, err := db.Query(models.NewQuery(query, nil))
	if err != nil {
		return nil, err
	}

	var subs []SubmissionDB
	if err := mapstructure.Decode(docs, &subs); err != nil {
		return nil, err
	}
	return subs, nil
}
