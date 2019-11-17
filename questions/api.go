// this file holds all of the question data, and hole id

package question

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Squwid/bytegolf/firestore"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

// Handler is the questions handler which takes care of all question related endpoint tasks
func Handler(w http.ResponseWriter, r *http.Request) {
	// stupid cors stuff
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS,GET")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodOptions {
		// this is for cors
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodGet {
		lightQs, err := listQuestions()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if len(lightQs) == 0 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("[]")) // write empty list if blank
			return
		}
		bs, err := json.Marshal(lightQs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(bs)
		return
	}
	// forbidden for other requests that are not cors, get
	w.WriteHeader(http.StatusForbidden)
}

// SingleHandler is the handler that you use to receive a single question from the database whether it is
// live or not. You should be able to get questions without being logged in
func SingleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS,GET")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodGet {
		qID := r.URL.Query().Get("hole") // get the hole from the query strings to query the hole
		if qID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var qs = []Question{}
		ctx := context.Background()
		iter := firestore.Client.Collection(collection).Where("ID", "==", qID).Documents(ctx)
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Errorf("error getting hole with id %s: %v", qID, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var q Question
			err = mapstructure.Decode(doc.Data(), &q)
			if err != nil {
				log.Errorf("error decoding hole with id %s: %v", qID, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			qs = append(qs, q)
		}
		log.Infof("expected 1 or 0 questions got %v", len(qs))

		if len(qs) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if len(qs) > 1 {
			log.Warnf("did not expect more than one question but got %v. Returning the first one", len(qs))
		}

		bs, err := json.Marshal(qs[0])
		if err != nil {
			log.Errorf("error marshalling hole %s: %v", qID, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(bs)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
}

// listQuestions gets a list of questions that have the Live bool
func listQuestions() ([]Light, error) {
	var qs = []Light{}
	ctx := context.Background()
	iter := firestore.Client.Collection(collection).Where("Live", "==", true).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var q Light
		err = mapstructure.Decode(doc.Data(), &q)
		if err != nil {
			log.Errorf("error decoding object: %v", err)
		} else {
			log.Debugf("got data back, parsing: %s", doc.Data())
			qs = append(qs, q)
		}
	}
	return qs, nil
}