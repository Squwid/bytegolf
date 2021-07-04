package jdoodle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Squwid/bytegolf/auth"
	"github.com/Squwid/bytegolf/compiler"
	"github.com/Squwid/bytegolf/db"
	"github.com/Squwid/bytegolf/globals"
	"github.com/Squwid/bytegolf/holes"
	"github.com/Squwid/bytegolf/models"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func SubmissionHandler(w http.ResponseWriter, r *http.Request) {
	holeID := mux.Vars(r)["hole"]
	log := logrus.WithFields(logrus.Fields{
		"Hole":   holeID,
		"Action": "NewSubmission",
		"IP":     r.RemoteAddr,
	})

	claims := auth.LoggedIn(r)
	if claims == nil {
		log.Infof("User not authenticated")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log = log.WithField("User", claims.BGID)

	// Make sure hole exists
	out, err := db.Get(models.NewGet(db.HoleCollection().Doc(holeID), nil))
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
	var hole holes.Hole
	if err := mapstructure.Decode(out, &hole); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithError(err).Errorf("Error parsing hole from db")
		return
	}
	if !hole.Active {
		w.WriteHeader(http.StatusNotFound)
		log.Warnf("got inactive hole")
		return
	}

	var userInput UserInput
	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		http.Error(w, "Bad Input: "+err.Error(), http.StatusBadRequest)
		log.WithError(err).Errorf("Error parsing user input")
		return
	}
	if valid, msg := userInput.validate(); !valid {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	log.Infof("Hole active & exists, getting test cases")

	tests, err := holes.GetTests(hole.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithError(err).Errorf("Error getting test cases")
		return
	}
	log = log.WithField("Tests", len(tests))

	// Wait for tests to complete
	var ch = make(chan compileResult, len(tests))
	for _, test := range tests {

		go func(test holes.Test, input UserInput) {
			compileInput := userInput.Input(test.Input)
			go compiler.Compile(compileInput)

			compileOutput := <-compileInput.response
			var result = compileResult{test: &test}

			if compileOutput.Err != nil {
				result.err = compileOutput.Err
				ch <- result

			} else if compileOutput.StatusCode != http.StatusOK {
				result.err = fmt.Errorf("got bad status code %v from compiler", compileOutput.StatusCode)
				ch <- result
			} else {
				var output Output
				if err := json.NewDecoder(compileOutput.Body).Decode(&output); err != nil {
					result.err = err
					ch <- result
					return
				}

				result.output = &output
				correct, err := test.Check(output.Output)
				if err != nil {
					result.err = err
					ch <- result
					return
				}
				result.correct = correct
				ch <- result
			}
		}(test, userInput)

	}

	var sub = compiler.NewSubmissionDB(holeID, globals.BGID, userInput.Script, userInput.Language, userInput.Version)
	i, correct, incorrect := 0, 0, 0
	timeout := time.NewTimer(time.Second * 15)

	// Wait for all tests to be done or timeout
	for {
		select {
		case out := <-ch:
			if out.err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.WithField("Test", out.test.ID).WithError(err).Errorf("Error compiling test")
				return
			}

			sub.AddTest(out.test.ID, out.correct)
			log.WithField("TestID", out.test.ID).WithField("Output", out.output.Output).Infof("Correct: %v", out.correct)

			i++
			if out.correct {
				correct++
			} else {
				incorrect++
			}

		case <-timeout.C:
			w.WriteHeader(http.StatusInternalServerError)
			log.Warnf("Compiler timed out")
			return
		}

		if i == len(tests) {
			break
		}
	}
	close(ch)

	if err := db.Store(sub); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithError(err).Errorf("Error storing submission")
		return
	}

	bs, _ := json.Marshal(struct {
		Correct bool  `json:"Correct"`
		Length  int64 `json:"Length"`
	}{
		Correct: incorrect == 0,
		Length:  sub.Length,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
	log.WithFields(logrus.Fields{
		"Correct":        incorrect == 0,
		"Length":         len(userInput.Script),
		"CorrectCount":   correct,
		"IncorrectCount": incorrect,
	}).Infof("Successful compile request")
}
