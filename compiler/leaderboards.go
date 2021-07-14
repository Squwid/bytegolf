package compiler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/Squwid/bytegolf/auth"
	"github.com/Squwid/bytegolf/db"
	"github.com/Squwid/bytegolf/models"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

type Entry struct {
	ID       string
	Language string
	Version  string
	Length   int64
	HoleID   string
	BGID     string

	// Following fields get populated by the BGID on a request basis
	GitName string
}

// GetUserFields populates the user fields based on BGID
func (e *Entry) GetUserFields() error {
	getter := models.NewGet(db.ProfileCollection().Doc(e.BGID), nil)
	profile, err := db.Get(getter)
	if err != nil {
		return err
	}

	e.GitName = profile["GithubUser"].(map[string]interface{})["Login"].(string)
	return nil
}

// LeaderboardQuery is a wrapper query function to list leaderboard submissions for a given hole.
// Since leaderboard have to be unique by BGID, that query is handled in this function. If the limit is not
// reached it will return a slice of the users it was able to get
func LeaderboardQuery(query firestore.Query, limit int) ([]Entry, error) {
	var users = make(map[string]bool) // Map users in O(1) time
	var entries = []Entry{}

	iter := query.Documents(context.Background())
	for {
		if len(entries) == limit {
			break
		}

		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, err
		}

		var sub SubmissionDB
		if err := doc.DataTo(&sub); err != nil {
			return nil, err
		}

		// User already exists in the leaderboard. LB has to be unique
		if users[sub.BGID] {
			continue
		}

		users[sub.BGID] = true
		entry := sub.Entry()
		entry.GetUserFields()
		entries = append(entries, entry)
	}

	return entries, nil
}

// Possible query strings:
//   "hole": ID of the hole *REQUIRED*
//   "limit": Limit of entry numbers. Max of 10. Defaults to 10.
//   "lang": Specific language to query for. Defaults to all languages.
func LeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	log := logrus.WithFields(logrus.Fields{
		"Action": "Leaderboard",
	})

	claims := auth.LoggedIn(r)
	if claims != nil {
		log = log.WithField("User", claims.BGID)
	}

	// holeID querystring is
	hole := r.URL.Query().Get("hole")
	if hole == "" {
		log.Warnf("Required querystring 'hole' missing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log = log.WithField("Hole", hole)

	query := db.SubmissionsCollection().Where("Correct", "==", true).Where("HoleID", "==", hole).
		OrderBy("Length", firestore.Asc).OrderBy("SubmittedTime", firestore.Asc)

	// Parse limit query string and make sure its valid 👍
	var limit = 5
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		tLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			log.WithError(err).Warnf("Invalid limit input '%s'", limitStr)
		} else if tLimit > 10 || limit <= 0 {
			log.Warnf("Invalid limit input '%v'", tLimit)
		} else {
			limit = tLimit
		}
	}
	log = log.WithField("Limit", limit)

	if lang := r.URL.Query().Get("lang"); lang != "" {
		log = log.WithField("Lang", lang)
		query = query.Where("Language", "==", lang)
	}

	leaders, err := LeaderboardQuery(query, limit)
	if err != nil {
		log.WithError(err).Errorf("Error querying for leaderboards")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(leaders)
	if err != nil {
		log.WithError(err).Errorf("Error marshalling leaders")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
	log.Infof("Got %v leaders", len(leaders))
}
