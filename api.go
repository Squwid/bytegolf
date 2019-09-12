package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Squwid/bytegolf/firestore"
	log "github.com/sirupsen/logrus"
)

// code is the entire request body that contains data about tests,
// input, language
type code struct {
	UID      string `json:"uid"`
	Language string `json:"language"`
	CodeBody string `json:"code_body"`
	Tests    []test `json:"tests"`
	Image    string `json:"image"`
}

type test struct {
	UID   string   `json:"uid"`
	Input []string `json:"input"`
}

type results struct {
	UID         string `json:"uid"`
	TestResults struct {
		UID      string    `json:"uid"`
		ExitCode int       `json:"exit_code"`
		Stdout   string    `json:"std_out"`
		Stderr   string    `json:"std_err"`
		Time     time.Time `json:"time"`
		Memory   string    `json:"memory"`
	}
}

func compile(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		log.Warnf("got %s request for /compile but expected 'POST'", req.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	start := time.Now()
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorf("error parsing request body for /compile: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	// i could unmarshal the code struct here and then remarshal it but
	// that seems useless (maybe for debugging at some point)
	creq, err := http.NewRequest(http.MethodPost, compilerURI, bytes.NewReader(bs))
	if err != nil {
		log.Errorf("error creating new request for /compile: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// this section code take up to a minute or so (3 min timeout)
	resp, err := http.DefaultClient.Do(creq)
	if err != nil {
		log.Errorf("error sending request to compile code: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if resp.StatusCode == http.StatusRequestTimeout {
		log.Errorf("request for /compile timed out...")
		w.WriteHeader(http.StatusRequestTimeout)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Errorf("got status code %v from /compile request", resp.Status)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bsr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("error reading response body from /compile request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	overall := time.Since(start)

	// this go function should upload this data to the database of all requests
	go func(req, resp []byte, dur time.Duration) {
		err := firestore.StoreData("compiles", map[string]interface{}{
			"code_request":     string(req),
			"code_results":     string(resp),
			"request_duration": overall.String(),
		})
		if err != nil {
			log.Errorf("Error storing fire data store: %v", err)
			return
		}
	}(bs, bsr, overall)
	log.Infof("overall time for compile request: %v", overall)
	w.Write(bsr)
}
