package auth

import (
	"github.com/Squwid/bytegolf/db"
	"github.com/Squwid/bytegolf/models"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

// Bytegolf returns a bytegolf user based on the github user. Will create one
// if one does not exist yet
func Bytegolf(ghu *GithubUser) (*BytegolfUser, error) {
	user, err := bytegolfUserFromGitID(ghu.ID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		// User does not exist, create one & store it
		logrus.WithField("Github", ghu.ID).Infof("Bytegolf user did not exist, created one")
		user = NewBytegolfUser(*ghu)
		if err := db.Store(user); err != nil {
			return nil, err
		}
	} else {
		// TODO: Check last updated time to see if it should be re-stored
		logrus.WithField("BGID", user.BGID).WithField("Github", ghu.ID).Infof("Found existing bytegolf user")
	}

	return user, nil
}

// bytegolfUserFromGitID will get a BytegolfUser from the table using a gitID, will return nil, nil if one
// does not exist
func bytegolfUserFromGitID(gitID int64) (*BytegolfUser, error) {
	query := db.ProfileCollection().Where("GithubUser.ID", "==", gitID).Limit(1)
	docs, err := db.Query(models.NewQuery(query, nil))
	if err != nil {
		return nil, err
	}

	// Parse user
	var bgu BytegolfUser
	if err := mapstructure.Decode(docs[0], &bgu); err != nil {
		return nil, err
	}

	return &bgu, nil
}
