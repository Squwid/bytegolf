package submissions

import (
	"context"

	"cloud.google.com/go/firestore"
	fs "github.com/Squwid/bytegolf/firestore"
	"github.com/Squwid/bytegolf/globals"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/iterator"
)

// GetBestSubmissionsOnHole gets all of the full submissions for a specific hole
func GetBestSubmissionsOnHole(hole string, max int) (FullSubmissions, error) {
	ctx := context.Background()
	iter := fs.Client.Collection(globals.SubmissionsTable()).Where("Short.Correct", "==", true).
		OrderBy("Short.Length", firestore.Asc).Documents(ctx)

	// contains is a function to see if the bgid is already in the slice
	contains := func(ss FullSubmissions, bgid string) bool {
		for _, s := range ss {
			if s.Short.BGID == bgid {
				return true
			}
		}
		return false
	}

	var subs = FullSubmissions{}
	for {
		// Break if the length is long enough
		if len(subs) >= max {
			break
		}

		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		// Parse document to submission
		var sub FullSubmission
		if err := mapstructure.Decode(doc.Data(), &sub); err != nil {
			return nil, err
		}

		// Only append if user's BGID is not already there
		if !contains(subs, sub.Short.BGID) {
			subs = append(subs, sub)
		}
	}

	return subs, nil
}

// GetSubmissionByID gets a full submission by using an ID
func GetSubmissionByID(id string) (*FullSubmission, error) {
	ctx := context.Background()
	iter := fs.Client.Collection(globals.SubmissionsTable()).Where("Short.UUID", "==", id).Limit(1).Documents(ctx)

	var sub FullSubmission

	// User iter without a loop
	doc, err := iter.Next()
	if err == iterator.Done {
		// Not found
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Decode
	if err := mapstructure.Decode(doc.Data(), &sub); err != nil {
		return nil, err
	}

	return &sub, nil
}

// Store stores a full submission in the submissions table
func (sub FullSubmission) Store() error {
	ctx := context.Background()

	_, _, err := fs.Client.Collection(globals.SubmissionsTable()).Add(ctx, sub)
	return err
}

// GetBestPlayerSubmission gets a best score on a specific hole
func GetBestPlayerSubmission(bgid, hole string) (*FullSubmission, error) {
	ctx := context.Background()
	iter := fs.Client.Collection(globals.SubmissionsTable()).Where("Short.BGID", "==", bgid).
		Where("Short.HoleID", "==", hole).Where("Short.Correct", "==", true).OrderBy("Short.Length", firestore.Asc).
		Limit(1).Documents(ctx)

	var sub FullSubmission

	// Iter but only expecting one
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Got the best submission
	if err := mapstructure.Decode(doc.Data(), &sub); err != nil {
		return nil, err
	}

	return &sub, nil
}
