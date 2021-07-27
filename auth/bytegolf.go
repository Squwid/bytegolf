package auth

import (
	"fmt"

	"github.com/Squwid/bytegolf/db"
	"github.com/Squwid/bytegolf/models"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BGUser returns a bytegolf user based on the github user. Will create one
// if one does not exist yet
func BGUser(ghu *GithubUser) (*BytegolfUser, error) {
	user, err := bgUserFromGit(ghu.ID)
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

// bgUserFromGit will get a BytegolfUser from the table using a gitID, will return nil, nil if one
// does not exist
func bgUserFromGit(gitID int64) (*BytegolfUser, error) {
	getter := models.NewGet(db.ProfileCollection().Doc(fmt.Sprintf("%v", gitID)), nil)
	doc, err := db.Get(getter)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, err
	}

	var bgu BytegolfUser
	if err := mapstructure.Decode(doc, &bgu); err != nil {
		return nil, err
	}

	return &bgu, nil
}
