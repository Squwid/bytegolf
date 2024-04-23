package bg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Squwid/bytegolf/auth"
	"github.com/Squwid/bytegolf/compiler"
	"github.com/Squwid/bytegolf/db"
	"github.com/Squwid/bytegolf/holes"
	"github.com/Squwid/bytegolf/models"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SubmissionHandler is the HTTP handler for the Bytegolf Compiler
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
	validationOutput := userInput.validate()
	if !validationOutput.valid {
		log.Infof("Invalid compile request: %s", validationOutput.msg)
		http.Error(w, validationOutput.msg, http.StatusBadRequest)
		return
	}
	log.Debugf("Hole active & exists, getting test cases")
	log.WithField("Language", userInput.Language).WithField("Version",
		userInput.Version).Debugf("Valid request")

	tests, err := holes.GetTests(hole.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithError(err).Errorf("Error getting test cases")
		return
	}
	log = log.WithField("Tests", len(tests))
	log.Debugf("Test cases: %+v", tests)

	// Run each test individually and compare results.
	var ch = make(chan compileResult, len(tests))
	for _, test := range tests {

		go func(test holes.Test, input UserInput) {
			compileInput := userInput.Input(test.Input, *validationOutput.language)
			go compiler.Compile(compileInput)

			compileOutput := <-compileInput.response
			var result = compileResult{test: &test}

			if compileOutput.Err != nil {
				result.err = compileOutput.Err
				ch <- result

			} else if compileOutput.StatusCode != http.StatusOK {
				result.err = fmt.Errorf("got bad status code %v from compiler",
					compileOutput.StatusCode)
				ch <- result
			} else {
				var outputs []Output
				if err := json.NewDecoder(compileOutput.Body).
					Decode(&outputs); err != nil {
					result.err = err
					ch <- result
					return
				}

				if len(outputs) != 1 {
					result.err = fmt.Errorf("expected 1 output, got %v", len(outputs))
					ch <- result
					return
				}

				result.output = &outputs[0]
				correct, err := test.Check(result.output.StdOut)
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

}
