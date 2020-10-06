package submissions

import (
	"context"

	"cloud.google.com/go/firestore"
	fs "github.com/Squwid/bytegolf/firestore"
	"github.com/Squwid/bytegolf/models"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const submissionCollection = "submissions"

func storeSubmissionDB(sub models.SubmissionDB) error {
	_, err := fs.Client.Collection(submissionCollection).Doc(sub.ID).Set(context.Background(), sub)
	return err
}

func getDBSubmission(id string) (*models.SubmissionDB, error) {
	doc, err := fs.Client.Collection(submissionCollection).Doc(id).Get(context.Background())
	if status.Code(err) == codes.NotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Parse the submission
	var sub models.SubmissionDB
	if err := doc.DataTo(&sub); err != nil {
		return nil, err
	}

	return &sub, nil
}

func getSingleBestSubmissionOnHole(holeID, bgid string) (*models.SubmissionDB, error) {
	ctx := context.Background()

	// Get players best score
	iter := fs.Client.Collection(submissionCollection).Where("MetaData.Correct", "==", true).
		Where("HoleID", "==", holeID).Where("BGID", "==", bgid).OrderBy("MetaData.Length", firestore.Asc).
		Limit(1).Documents(ctx)

	docs, err := iter.GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, nil
	}

	var sub models.SubmissionDB
	if err := docs[0].DataTo(&sub); err != nil {
		return nil, err
	}

	return &sub, nil
}

// getBestSubmissionsOnHole gets the best submissions on a hole, but one per bgid
func getBestSubmissionsOnHole(holeID string, max int) ([]models.SubmissionDB, error) {
	ctx := context.Background()

	// Get correct answers, 1 per bg for a single hole
	iter := fs.Client.Collection(submissionCollection).Where("MetaData.Correct", "==", true).
		Where("HoleID", "==", holeID).OrderBy("MetaData.Length", firestore.Asc).Documents(ctx)

	contains := func(subs []models.SubmissionDB, bgid string) bool {
		for _, sub := range subs {
			if sub.BGID == bgid {
				return true
			}
		}
		return false
	}

	var subs = []models.SubmissionDB{}
	for {
		// Break if length is long enough
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

		// Parse document to SubmissionDB object
		var sub models.SubmissionDB
		if err := doc.DataTo(&sub); err != nil {
			return nil, err
		}

		// Only append if user's BGID not already there
		if !contains(subs, sub.BGID) {
			subs = append(subs, sub)
		}
	}

	return subs, nil
}

// getUserPastSubmissions gets a users past submissions, if holeID is left blank, then gets all holes
func getUserPastSubmissions(bgid, holeID string, max int) ([]models.SubmissionDB, error) {
	ctx := context.Background()

	// Get all answers, correct or not, sorting by created at
	query := fs.Client.Collection(submissionCollection).Where("BGID", "==", bgid)

	// have optional query holeID
	if holeID != "" {
		query = query.Where("HoleID", "==", holeID)
	}

	iter := query.OrderBy("CreatedAt", firestore.Asc).Documents(ctx)

	var subs = []models.SubmissionDB{}
	for {
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

		var sub models.SubmissionDB
		if err := doc.DataTo(&sub); err != nil {
			return nil, err
		}

		subs = append(subs, sub)
	}

	return subs, nil
}
