package holes

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func ListHoles(w http.ResponseWriter, r *http.Request) {
	log := logrus.WithFields(logrus.Fields{
		"Action": "ListHoles",
	})

	holes, err := listHoles()
	if err != nil {
		log.WithError(err).Errorf("Error listing holes")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(holes)
	if err != nil {
		log.WithError(err).Errorf("error marshalling holes")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}
