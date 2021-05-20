package holes

import (
	"encoding/json"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/Squwid/bytegolf/db"
	"github.com/Squwid/bytegolf/models"
	"github.com/sirupsen/logrus"
)

func ListHoles(w http.ResponseWriter, r *http.Request) {
	log := logrus.WithFields(logrus.Fields{
		"Action": "ListHoles",
	})

	query := db.HoleCollection().OrderBy("CreatedAt", firestore.Desc).Where("Active", "==", true).Limit(100)

	hs, err := db.Query(models.NewQuery(query, transformHole))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.WithError(err).Errorf("Error querying holes")
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
