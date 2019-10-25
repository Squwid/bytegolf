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
