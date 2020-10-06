package submissions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Squwid/bytegolf/auth"
	"github.com/Squwid/bytegolf/holes"
	"github.com/Squwid/bytegolf/jdoodle"
	"github.com/Squwid/bytegolf/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// GetBestSubmissionOnHoleHandler gets the best submission on a hole
func GetBestSubmissionOnHoleHandler(w http.ResponseWriter, r *http.Request) {
	hole := mux.Vars(r)["hole"]

	log := logrus.WithFields(logrus.Fields{
		"Action": "BestSubmissionOnHole",
		"IP":     r.RemoteAddr,
		"Hole":   hole,
	})

	loggedIn, claims := auth.LoggedIn(r)
	if !loggedIn {
		log.Warnf("User not logged in")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	log = log.WithField("BGID", claims.BGID)

	dbSub, err := getSingleBestSubmissionOnHole(hole, claims.BGID)
	if err != nil {
		log.WithError(err).Errorf("Error getting best submission for hole")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if dbSub == nil {
		http.Error(w, "not found", http.StatusNotFound)
		log.Warnf("No best submission")
		return
	}

	sub := dbSub.Frontend()

	bs, err := json.Marshal(sub)
	if err != nil {
		log.WithError(err).Errorf("Error marshalling submission")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	log.Infof("Got best submission")

	w.Write(bs)
}

// GetBestHoleSubmissionsHandler gets the best correct submissions for a specific hole
func GetBestHoleSubmissionsHandler(w http.ResponseWriter, r *http.Request) {
	const maxSubs = 5

	hole := mux.Vars(r)["hole"]

	log := logrus.WithFields(logrus.Fields{
		"Action": "BestHoleSubmission",
		"IP":     r.RemoteAddr,
		"Hole":   hole,
	})

	// Get best hole subs
	dbSubs, err := getBestSubmissionsOnHole(hole, maxSubs)
	if err != nil {
		log.WithError(err).Errorf("Error getting best subs")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	subs := models.SubmissionTransform(dbSubs)

	bs, err := json.Marshal(subs)
	if err != nil {
		log.WithError(err).Errorf("Error marshalling submissions")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Count", fmt.Sprintf("%v", len(subs)))
	w.Header().Set("Content-Type", "application/json")

	log.Infof("Got %v submissions", len(subs))

	w.Write(bs)
}

// GetMySubmissionsHandler gets a logged in user's past 20 past submissions
func GetMySubmissionsHandler(w http.ResponseWriter, r *http.Request) {
	const maxSubs = 20

	log := logrus.WithFields(logrus.Fields{
		"Action": "GetMySubs",
		"IP":     r.RemoteAddr,
	})

	// Optional query string
	holeID := r.URL.Query().Get("hole")
	if holeID != "" {
		log = log.WithField("Hole", holeID)
	}

	loggedIn, claims := auth.LoggedIn(r)
	if !loggedIn {
		log.Warnf("User not logged in")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	log = log.WithField("BGID", claims.BGID)

	// Get the user's past submissions
	dbSubs, err := getUserPastSubmissions(claims.BGID, holeID, maxSubs)
	if err != nil {
		log.WithError(err).Errorf("Error getting past submissions")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Transform from database objects
	subs := models.SubmissionTransform(dbSubs)

	bs, err := json.Marshal(subs)
	if err != nil {
		log.WithError(err).Errorf("Error marshalling submissions")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Count", fmt.Sprintf("%v", len(subs)))
	w.Header().Set("Content-Type", "application/json")

	log.Infof("Got %v submissions", len(subs))

	w.Write(bs)
}

// NewSubmissionHandler is the handler to submit a new submission
func NewSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	hole := mux.Vars(r)["hole"]

	log := logrus.WithFields(logrus.Fields{
		"Action": "NewSubmission",
		"Hole":   hole,
		"IP":     r.RemoteAddr,
	})

	loggedIn, claims := auth.LoggedIn(r)
	if !loggedIn {
		log.Warnf("User not logged in")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	log = log.WithFields(logrus.Fields{
		"BGID": claims.BGID,
	})

	// Get DBHole Object
	holeDB, err := holes.GetDBHole(hole)
	if err != nil {
		log.WithError(err).Errorf("Error getting hole")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Hole doesnt exist, return a 404
	if holeDB == nil {
		log.Warnf("Hole does not exist")
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	log.Infof("New incoming submission")

	// Parse input
	var submission models.IncomingSubmission
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		log.WithError(err).Errorf("Error unmarshalling body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate input
	if submission.Code == "" || submission.Language == "" {
		log.Warnf("Invalid input")
		http.Error(w, "code and lanugage cannot be blank", http.StatusBadRequest)
		return
	}

	// Get the language and version index for JDoodle
	language, versionIndex, err := jdoodle.MapLanguage(submission.Language)
	if err != nil {
		// Language was not found
		log.WithError(err).Errorf("Language not found")
		http.Error(w, "language not found", http.StatusBadRequest)
		return
	}

	var compileChan = make(chan struct {
		CompileIn  models.CompileInput
		CompileOut models.CompileOutput

		TestIn  models.TestCaseInput
		TestOut models.TestCaseOutput

		err error
	}, len(holeDB.Tests))

	var wg = sync.WaitGroup{}
	wg.Add(len(holeDB.Tests))

	// Compile each test case
	for _, test := range holeDB.Tests {
		go func(test models.TestCaseInput) {
			var compile = struct {
				CompileIn  models.CompileInput
				CompileOut models.CompileOutput

				TestIn  models.TestCaseInput
				TestOut models.TestCaseOutput

				err error
			}{
				TestIn: test,
			}

			compile.CompileIn = models.CompileInput{
				Code:         submission.Code,
				StdIn:        test.Input,
				Language:     language,
				VersionIndex: versionIndex,
			}

			// Send request to JDoodle
			compileOutput, err := jdoodle.SendJdoodle(compile.CompileIn)
			if err != nil {
				compile.err = err
				compileChan <- compile
				return
			}

			compile.CompileOut = *compileOutput

			// Create a TestOut
			testOutput, err := jdoodle.CheckOutput(test, compile.CompileOut)
			if err != nil {
				compile.err = err
				compileChan <- compile
				return
			}

			compile.TestOut = *testOutput

			compileChan <- compile
			wg.Done()
		}(test)
	}

	// Wait for all compiles to be done
	// TODO: Timeout
	wg.Wait()
	close(compileChan)

	var submissionDB = models.SubmissionDB{
		ID:        uuid.New().String(),
		HoleID:    hole,
		BGID:      claims.BGID,
		CreatedAt: time.Now().UTC(),
		Jdoodles:  []models.Jdoodle{},
		Tests:     []models.SubmissionDBTest{},
		MetaData:  models.SubmissionMetaData{},
	}

	// Iterate over each test case
	for compile := range compileChan {
		if compile.err != nil {
			// TODO: What are the possible errors here?
			log.WithError(err).Errorf("Error compiling test %s", compile.TestIn.ID)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Get rid of the code section, since its not needed in the database
		compile.CompileIn.Code = ""

		submissionDB.Jdoodles = append(submissionDB.Jdoodles, models.Jdoodle{
			CompileInput:  compile.CompileIn,
			CompileOutput: compile.CompileOut,
		})

		// Create a test object with TestIn and TestOut
		fullTest := models.SubmissionDBTest{
			TestInput:  compile.TestIn,
			TestOutput: compile.TestOut,
		}

		submissionDB.Tests = append(submissionDB.Tests, fullTest)
	}

	// Add metadata
	submissionDB.MetaData.Code = submission.Code
	submissionDB.MetaData.Length = len(submission.Code)
	submissionDB.MetaData.Language = language

	// Check if all tests pass
	var passedAll = true
	for _, test := range submissionDB.Tests {
		if !test.TestOutput.Correct {
			passedAll = false
		}
	}
	submissionDB.MetaData.Correct = passedAll

	// Store submission to DB
	if err := storeSubmissionDB(submissionDB); err != nil {
		log.WithError(err).Errorf("Error storing submission")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Infof("Successfully stored submission %s to the database", submissionDB.ID)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"success": "stored %s"}`, submissionDB.ID)))
}
