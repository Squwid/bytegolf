package leaderboard

import (
	"encoding/json"
	"net/http"

	"github.com/Squwid/bytegolf/compiler"
	"github.com/Squwid/bytegolf/github"
	"github.com/Squwid/bytegolf/question"
	"github.com/Squwid/bytegolf/sess"
	log "github.com/sirupsen/logrus"
)

// Spot represents a single leaderboard spot
type Spot struct {
	compiler.ShortSubmission

	GithubURI string `json:"github_uri"`
	Username  string `json:"username"`
}

// Handler is the handler for the leaderboard
func Handler(w http.ResponseWriter, r *http.Request) {
	// the user doesnt have to be logged in to see the leaderboard scores
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodOptions {
		// this is for cors
		w.WriteHeader(http.StatusOK)
		return
	}

	// method can only be get
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	hole := r.URL.Query().Get("hole")
	if hole == "" {
		log.Warnf("Request to list leaderboard spots but 'hole' query string was missing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// i only allow for 3 because i get charged per read and if the user lists -1 or something
	// stupid i would not be happy
	const max = 3

	q, err := question.GetQuestion(hole)
	if err != nil {
		log.Errorf("Error checking if question %s exists: %v", hole, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if q == nil {
		// question was not found so return a 404 with a desc that says youre lost (uzi)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "hole not found"}`))
		return
	}

	// get the best score for a user with their github + stuff
	best := r.URL.Query().Get("best")
	if best != "" {
		loggedIn, s, err := sess.LoggedIn(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !loggedIn || s == nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		sub, err := compiler.GetBestSubmission(s.BGID, hole)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if sub == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		ss := sub.TransformToShort()
		user, err := github.RetreiveUser(sub.BGID)
		if err != nil || user == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		spot := Spot{
			ShortSubmission: ss,
			GithubURI:       user.GithubURI,
			Username:        user.Username,
		}
		bs, err := json.Marshal(spot)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(bs)
		return
	}

	// the hole was found so lets get the leaders!
	fullSubs, err := compiler.GetBestSubmissionsOnHole(q.ID, max)
	if err != nil {
		log.Errorf("Error getting best submissions on %s: %v", hole, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// this makes it way easier than marshalling
	if len(fullSubs) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	// we have full submissions but they need to get changed to the leaderboard `Spot` struct
	var spots = []Spot{}

	for _, sub := range fullSubs {
		var spot = Spot{
			ShortSubmission: compiler.ShortSubmission{
				ID:            sub.UUID,
				Correct:       sub.Correct,
				Language:      sub.Exe.Language,
				Score:         sub.Length,
				SubmittedTime: sub.SubmittedTime,
			},
		}

		user, err := github.RetreiveUser(sub.BGID)
		if err != nil || user == nil {
			// the user was not found somehow (this should honestly never happen) OR
			// there was an unexpected error
			// so put my github url in there
			spot.GithubURI = "https://github.com/Squwid"
			spot.Username = "Not Found"
		} else {
			spot.GithubURI = user.GithubURI
			spot.Username = user.Username
		}
		spots = append(spots, spot)
	}

	bs, err := json.Marshal(spots)
	if err != nil {
		log.Errorf("Error marshalling %v leaderboard spots for hole %s: %v", max, hole, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}
