package compiler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Squwid/bytegolf/firestore"
	"github.com/Squwid/bytegolf/question"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// Possible errors that can come from jdoodle
var (
	ErrOutOfCredits = errors.New("Out of jdoodle api credits for this account")
	ErrGotBadStatus = errors.New("got a bad status back from the server with no response")
)

// RunTests will take a hole and an execute and run the tests, it will return a
// FullSubmission
func (exe Execute) RunTests(q question.Question, bgid string) (*FullSubmission, error) {
	// need to run every test with the same execute code but different std in

	uid := uuid.New().String()
	var submission = FullSubmission{
		UUID:          uid,
		BGID:          bgid,
		HoleID:        q.ID,
		Length:        len(exe.Script),
		SubmittedTime: time.Now(),
		Exe:           exe,
	}

	// ch is the channel that the outputs get sent through, it will wait for them
	// all to be done
	type TestWithErrors struct {
		Output TestOutput
		Err    error
	}
	var ch = make(chan TestWithErrors)

	var tests int
	for _, tc := range q.TestCases {
		tests++
		// each go func should send the specific test to jdoodle and return a TestOutput through
		// the channel that will be waiting
		go func(tcase question.Test) {
			// TODO: better logs with ids because it would be impossible to know which is which
			bs, err := json.Marshal(struct {
				Execute
				StdIn        string `json:"stdin"`
				ClientID     string `json:"clientId"`
				ClientSecret string `json:"clientSecret"`
			}{
				Execute:      exe,
				StdIn:        tcase.Input,
				ClientID:     jdoodleClient.Client,
				ClientSecret: jdoodleClient.Secret,
			})

			log.Infof("Request: %v", string(bs))
			if err != nil {
				log.Errorf("Error parsing execute: %v", err)
				ch <- TestWithErrors{Err: err} // pass the error through the channel for that reply
				return
			}

			// make the request to send through
			jreq, err := http.NewRequest(http.MethodPost, compileURI, bytes.NewReader(bs))
			if err != nil {
				log.Errorf("Error creating jdoodle request: %v", err)
				ch <- TestWithErrors{Err: err}
				return
			}

			jreq.Header.Set("Content-Type", "application/json")
			jresp, err := http.DefaultClient.Do(jreq)
			if err != nil {
				log.Errorf("Error sending request to jdoodle: %v", err)
				ch <- TestWithErrors{Err: err}
				return
			}

			// check the status codes to see if jdoodle says were out of tokens or other errors
			if jresp.StatusCode == 409 {
				log.Warnf("Ran out of credits!")
				ch <- TestWithErrors{Err: ErrOutOfCredits}
				return
			} else if jresp.StatusCode != 200 {
				log.Warnf("Got bad status code of %v from jdoodle: %v", jresp.StatusCode)
				ch <- TestWithErrors{Err: ErrGotBadStatus}
				return
			}

			// read the entire response at once
			bs, err = ioutil.ReadAll(jresp.Body)
			if err != nil {
				log.Errorf("Error reading output body: %v", err)
				ch <- TestWithErrors{Err: err}
				return
			}
			defer jresp.Body.Close()

			// if the status was 200 parse the output into a ExecuteResponse struct and check
			// to see that it was correct
			var er ExecuteResponse
			err = json.Unmarshal(bs, &er)
			if err != nil {
				log.Errorf("Error parsing jdoodle response: %v", err)
				ch <- TestWithErrors{Err: err}
				return
			}

			// check to see if the uesr got the question right. for now it only checks
			// output directly to output, but maybe in the future do regex
			var correct = strings.TrimSpace(er.Output) == strings.TrimSpace(tcase.ExpectedOutput)

			// create the full object with no errors and pass it through
			var to = TestOutput{
				tcase,
				er,
				correct,
			}

			ch <- TestWithErrors{Output: to}
			log.Infof("Successfully compiled and tested")
			return
		}(tc)
	}

	// no tests were ran so just return what we had already
	if tests == 0 {
		log.Warnf("No tests for request %s", uid)
		err := firestore.StoreData("executes", uid, submission)
		if err != nil {
			log.Errorf("Error storing %s: %v", uid, err)
		}
		return &submission, err
	}

	var tos = []TestOutput{}

	// make a new structure that holds the execute and the client id and secret for jdoodle
	var c int // counter to know when to exit loop
	var overallCorrect = true
	for test := range ch {

		if test.Err != nil {
			return nil, test.Err
		}

		// the answer was wrong so make the overall correct to be false
		if !test.Output.Correct {
			overallCorrect = false
		}

		// the error was not null so append it to the test outputs
		tos = append(tos, test.Output)
		c++
		if c == tests {
			close(ch)
			break
		}
	}

	// set the rest of the submission stuff
	submission.TestOutputs = tos
	submission.Correct = overallCorrect

	log.Infof("Successfully ran %v tests for %s (%v), now storing...", len(tos), bgid, uid)
	err := firestore.StoreData("executes", uid, submission)
	if err != nil {
		log.Errorf("Error storing %s: %v", uid, err)
	}

	return &submission, err
}
