package leaderboard

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/Squwid/bytegolf/compiler"
	fs "github.com/Squwid/bytegolf/firestore"
	"github.com/Squwid/bytegolf/github"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

// Spot represents a single leaderboard spot
type Spot struct {
	compiler.ShortSubmission

	GithubURI string `json:"github_uri"`
	Username  string `json:"username"`

	bgid string // unexported but need this to look for dupes
}

func selectTopScores(hole string, max int) ([]Spot, error) {
	/*
		select top scores returns the top scorers for a specific hole.
		1. It needs to make sure that the BGID for each response is unique
		2. after getting the max amount of spots, map each bgid to player Username + github url
	*/

	ctx := context.Background()
	iter := fs.Client.Collection("executes").Where("HoleID", "==", hole).Where("Correct", "==", true).
		OrderBy("Length", firestore.Asc).Documents(ctx)

	var spots = []Spot{}

	contains := func(ss []Spot, bgid string) bool {
		for _, s := range ss {
			if s.bgid == bgid {
				return true
			}
		}
		return false
	}

	for {
		// if the length of spots is already maxed out, return
		if len(spots) >= max {
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

		var ts compiler.TotalStore
		err = mapstructure.Decode(doc.Data(), &ts)
		if err != nil {
			log.Errorf("Error decoding leaderboard submission for hole (%s): %v", hole, err)
			return nil, err
		}

		// only append to the slice if the user's bgid is not already there
		if !contains(spots, ts.BGID) {
			spots = append(spots, Spot{
				ShortSubmission: compiler.ShortSubmission{
					ID:            ts.UUID,
					Correct:       ts.Correct,
					Language:      ts.Exe.Language,
					Score:         ts.Length,
					SubmittedTime: ts.SubmittedTime,
				},
				bgid: ts.BGID,
			})
		}

	}

	// now we have the submissions, ,map the users bgid to an actual user
	// i could do some fancy stuff with go routines, but is it really worth it? maybe sometime down the road
	for i := range spots {
		user, err := github.RetreiveUser(spots[i].bgid)
		if err != nil && err == github.ErrNotFound {
			// if the user doesnt exist for some very odd reason, im going to make the name not found and link
			// to my own github
			log.Warnf("Submission %s could not find a matching user %s", spots[i].ID)
			spots[i].GithubURI = "https://github.com/Squwid"
			spots[i].Username = "Not Found"
			continue
		}
		if err != nil || user == nil {
			// make sure the user isnt null? i dont think this would ever happen but just incase
			log.Errorf("Error getting user %s from users table: %v", spots[i].bgid, err)
			spots[i].GithubURI = "https://github.com/Squwid"
			spots[i].Username = "Not Found"
			continue
		}

		spots[i].GithubURI = user.GithubURI
		spots[i].Username = user.Username
	}

	log.Infof("Request to get leaders for hole %s returned %s lb spots", hole, len(spots))
	return spots, nil
}
