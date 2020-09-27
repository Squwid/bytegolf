package auth

import (
	"context"
	"strings"
	"time"

	fs "github.com/Squwid/bytegolf/firestore"
	"github.com/Squwid/bytegolf/models"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

// Store stores a bytegolf user into the table, rewriting the entire old object
func Store(bgu *models.BytegolfUser) error {
	bgu.LastUpdatedTime = time.Now().UTC() // set last updated time

	ctx := context.Background()
	_, err := fs.Client.Collection(profileCollection).Doc(bgu.BGID).Set(ctx, *bgu)
	return err
}

// Bytegolf returns a bytegolf user based on the github user. Will create one
// if one does not exist yet
func Bytegolf(ghu *models.GithubUser) (*models.BytegolfUser, error) {
	user, err := bytegolfUserFromGitID(ghu.ID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		// User does not exist, create one & store it
		log.WithField("Github", ghu.ID).Infof("Bytegolf user did not exist, created one")
		user = NewBytegolfUser(*ghu)
		Store(user)
	} else {
		// TODO: Check last updated time to see if it should be re-stored
		log.WithField("BGID", user.BGID).WithField("Github", ghu.ID).Infof("Found existing bytegolf user")
	}

	return user, nil
}

// bytegolfUserFromGitID will get a BytegolfUser from the table using a gitID, will return nil, nil if one
// does not exist
func bytegolfUserFromGitID(gitID int64) (*models.BytegolfUser, error) {
	ctx := context.Background()
	docs, err := fs.Client.Collection(profileCollection).Where("GithubUser.ID", "==", gitID).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	// Only expectine one BytegolfUser per gitID
	if len(docs) > 1 {
		log.Warnf("Got %v BytegolfUsers but only expected 1", len(docs))
	} else if len(docs) == 0 {
		return nil, nil
	}

	// Parse user
	var bgu models.BytegolfUser
	if err := mapstructure.Decode(docs[0].Data(), &bgu); err != nil {
		return nil, err
	}

	return &bgu, nil
}

// bytegolfUserFromBGID gets a BytegolfUser from the table using a BytegolfID instead of a gitID.
// It will return nil, nil if one does not exist
func bytegolfUserFromBGID(bgid string) (*models.BytegolfUser, error) {
	ctx := context.Background()
	doc, err := fs.Client.Collection(profileCollection).Doc(bgid).Get(ctx)
	if err != nil && strings.Contains(err.Error(), "NotFound") {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Got document, parse it
	var bgu models.BytegolfUser
	if err := mapstructure.Decode(doc.Data(), &bgu); err != nil {
		return nil, err
	}

	return &bgu, nil
}
