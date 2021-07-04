package profiles

import (
	"encoding/json"
	"net/http"

	"github.com/Squwid/bytegolf/auth"
	"github.com/Squwid/bytegolf/db"
	"github.com/Squwid/bytegolf/models"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log := logrus.WithFields(logrus.Fields{
		"ID":     id,
		"Action": "GetProfile",
		"IP":     r.RemoteAddr,
	})

	claims := auth.LoggedIn(r)
	if claims != nil {
		log = log.WithField("User", claims.BGID)
	}

	getter := models.NewGet(db.ProfileCollection().Doc(id), transform)
	profile, err := db.Get(getter)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			w.WriteHeader(http.StatusNotFound)
			log.Warnf("Profile not found")
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		log.WithError(err).Errorf("Error getting profile")
		return
	}

	bs, err := json.Marshal(profile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithError(err).Error("Error unmarshalling user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
	log.Info("Got profile")
}

func transform(profile map[string]interface{}) error {
	delete(profile, "CreatedTime")
	delete(profile, "LastUpdatedTime")
	delete(profile["GithubUser"].(map[string]interface{}), "UpdatedAt")
	return nil
}
