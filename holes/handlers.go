package holes

import (
	"encoding/json"
	"net/http"

	"github.com/Squwid/bytegolf/auth"
	"github.com/Squwid/bytegolf/db"
	"github.com/Squwid/bytegolf/models"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListHoles lists 100 active holes
func ListHoles(w http.ResponseWriter, r *http.Request) {
	log := logrus.WithFields(logrus.Fields{
		"Action": "ListHoles",
	})

	claims := auth.LoggedIn(r)
	if claims != nil {
		log = log.WithField("User", claims.BGID)
	}

	query := db.HoleCollection().OrderBy("CreatedAt", firestore.Desc).
		Where("Active", "==", true).Limit(100)

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
	log.Infof("Listed %v holes", len(hs))
}

// GetHole gets a hole using an id
func GetHole(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log := logrus.WithFields(logrus.Fields{
		"ID":     id,
		"Action": "GetHole",
	})

	claims := auth.LoggedIn(r)
	if claims != nil {
		log = log.WithField("User", claims.BGID)
	}

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

	// TODO: Allow certain permissions to get holes that are not active
	if !hole["Active"].(bool) {
		w.WriteHeader(http.StatusNotFound)
		log.Warnf("Hole is inactive")
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

// GetTest gets a test using a hole and test id
func GetTest(w http.ResponseWriter, r *http.Request) {
	hole := mux.Vars(r)["hole"]
	id := mux.Vars(r)["id"]
	log := logrus.WithFields(logrus.Fields{
		"Test":   id,
		"Hole":   hole,
		"Action": "GetTest",
	})

	claims := auth.LoggedIn(r)
	if claims != nil {
		log = log.WithField("User", claims.BGID)
	}

	getter := models.NewGet(db.TestSubCollection(hole).Doc(id), nil)
	test, err := db.Get(getter)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			w.WriteHeader(http.StatusNotFound)
			log.Warnf("Test not found")
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		log.WithError(err).Errorf("Error getting test")
		return
	}

	var t Test
	if err := mapstructure.Decode(test, &t); err != nil {
		log.WithError(err).Errorf("Error decoding test case")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: Allow certain permissions to get holes that are not active
	if !t.Active {
		w.WriteHeader(http.StatusNotFound)
		log.Warnf("Test is inactive")
		return
	}

	bs, err := json.Marshal(t.ShortTest())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithError(err).Errorf("Error marshalling test")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
	log.Infof("Got test")
}

// ListTests gets all tests for a specific hole
func ListTests(w http.ResponseWriter, r *http.Request) {
	hole := mux.Vars(r)["hole"]
	log := logrus.WithFields(logrus.Fields{
		"Hole":   hole,
		"Action": "ListTests",
	})

	claims := auth.LoggedIn(r)
	if claims != nil {
		log = log.WithField("User", claims.BGID)
	}

	tests, err := GetTests(hole)
	if err != nil {
		log.WithError(err).Errorf("Error getting tests")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(tests.ShortTests())
	if err != nil {
		log.WithError(err).Errorf("Error marshalling tests")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
	log.Infof("Got %v test cases", len(tests))
}
