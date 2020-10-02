package submissions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/Squwid/bytegolf/auth"
	fs "github.com/Squwid/bytegolf/firestore"
	"github.com/Squwid/bytegolf/holes"
	"github.com/Squwid/bytegolf/jdoodle"
	"github.com/Squwid/bytegolf/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

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
		ID:       uuid.New().String(),
		HoleID:   hole,
		BGID:     claims.BGID,
		Jdoodles: []models.Jdoodle{},
		// Tests: models.SubmissionDBTest{
		// 	TestInputs:  []models.TestCaseInput{},
		// 	TestOutputs: []models.TestCaseOutput{},
		// },
		Tests:    []models.SubmissionDBTest{},
		MetaData: models.SubmissionMetaData{},
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

func storeSubmissionDB(sub models.SubmissionDB) error {
	const submissionsCollections = "submissions"
	_, err := fs.Client.Collection(submissionsCollections).Doc(sub.ID).Set(context.Background(), sub)
	return err
}
