package jdoodle

import (
	"encoding/json"
	"net/http"

	"github.com/Squwid/bytegolf/compiler"
	"github.com/Squwid/bytegolf/db"
	"github.com/Squwid/bytegolf/models"
	"github.com/gorilla/mux"
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

	// Make sure hole exists
	_, err := db.Get(models.NewGet(db.HoleCollection().Doc(holeID), nil))
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

	// TODO: Make sure hole is active

	log.Infof("Hole exists, getting test cases")

	// TODO: Get test cases

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

	// TODO: The code below all has to be changed to look at test output but ill worry about that later
	compileInput := userInput.Input("")
	go compiler.Compile(compileInput)

	compileOutput := <-compileInput.response
	if compileOutput.Err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithError(compileOutput.Err).Errorf("Error compiling")
		return
	}

	log.Infof("Sucessful compile request")

	bs, _ := json.Marshal(compileOutput.Out)
	w.Write(bs)
}
