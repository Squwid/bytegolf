package holes

import (
	"encoding/json"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/Squwid/bytegolf/db"
	"github.com/Squwid/bytegolf/models"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ListHoles(w http.ResponseWriter, r *http.Request) {
	log := logrus.WithFields(logrus.Fields{
		"Action": "ListHoles",
	})

	query := db.HoleCollection().OrderBy("CreatedAt", firestore.Desc).Where("Active", "==", true).Limit(100)

	hs, err := db.Query(models.NewQuery(query, transformHole))
	if err != nil {
		log.WithError(err).Errorf("Error querying holes")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(hs)
	if err != nil {
		log.WithError(err).Errorf("error marshalling holes")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}

func GetHole(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	log := logrus.WithFields(logrus.Fields{
		"ID":     id,
		"Action": "GetHole",
		"IP":     r.RemoteAddr,
	})

	getter := models.NewGet(db.HoleCollection().Doc(id), transformHole)
	hole, err := db.Get(getter)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			w.WriteHeader(http.StatusNotFound)
			log.Warnf("Hole not found")
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		log.WithError(err).Errorf("Error getting hole")
		return
	}

	bs, err := json.Marshal(hole)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithError(err).Errorf("Error unmarshalling hole")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
	log.Info("Got hole")
}
