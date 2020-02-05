// this file holds all of the question data, and hole id

package question

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"cloud.google.com/go/firestore"
	fs "github.com/Squwid/bytegolf/firestore"
	"github.com/Squwid/bytegolf/sess"
	"github.com/Squwid/bytegolf/util"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

// ListQuestionsHandler is the questions handler which takes care of all question related endpoint tasks
func ListQuestionsHandler(w http.ResponseWriter, r *http.Request) {
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
		var onlyLive = true
		if r.URL.Query().Get("allQs") == "true" {
			onlyLive = false
		}

		lightQs, err := listQuestions(onlyLive)
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

		var qs = []Light{}
		ctx := context.Background()
		iter := fs.Client.Collection(collection).Where("ID", "==", qID).Documents(ctx)
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
			var q Light
			err = mapstructure.Decode(doc.Data(), &q)
			if err != nil {
				log.Errorf("error decoding hole with id %s: %v", qID, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			qs = append(qs, q)
		}
		log.Debugf("expected 1 or 0 questions got %v", len(qs))

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
	if r.Method == http.MethodPost {
		// if the request is to create a new question then this function will
		// handle everything including permissions`
		CreateHole(w, r)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
}

// CreateHole is what gets called to create a new question, it will be denied if the uesr
// is not a game master
func CreateHole(w http.ResponseWriter, r *http.Request) {
	loggedIn, sess, err := sess.LoggedIn(r)
	if !loggedIn {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err != nil {
		log.Errorf("Error getting a session to create hole: %v", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if !sess.IsGamemaster() {
		log.Warnf("User %s tried to create a hole but insufficient permissions", sess.BGID)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Could not read new hole creation by %s: %v", sess.BGID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var q Question
	err = json.Unmarshal(bs, &q)
	if err != nil {
		log.Errorf("Error unmarshalling new hole by %s: %v", sess.BGID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// create will generate an id for the hole as well as storing it
	if err = q.create(); err != nil {
		log.Errorf("Error creating new question by %s: %v", sess.BGID, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(util.WrapError(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": "success"}`))
}

// listQuestions gets a list of questions that have the Live bool
func listQuestions(onlyLive bool) ([]Light, error) {
	var qs = []Light{}
	ctx := context.Background()

	// if onlyLive return only live questions, if not return every question
	var iter *firestore.DocumentIterator
	if onlyLive {
		// only live questions
		iter = fs.Client.Collection(collection).Where("Live", "==", true).Documents(ctx)
	} else {
		// all questions
		iter = fs.Client.Collection(collection).Documents(ctx)
	}

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
