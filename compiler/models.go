package compiler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

const compileURI = "https://api.jdoodle.com/v1/execute"

type Execute struct {
	Script       string `json:"script"`
	Language     string `json:"language"`
	VersionIndex string `json:"versionIndex"`

	// TODO: remove these from here and put them somewhere else
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type Response struct {
	Output     string `json:"output"`
	StatusCode int    `json:"statusCode"`
	Memory     string `json:"memory"`
	CPUTime    string `json:"cpuTime"`
}

// Possible errors that can come from jdoodle
var (
	ErrOutOfCredits = errors.New("Out of jdoodle api credits for this account")
	ErrGotBadStatus = errors.New("got a bad status back from the server with no response")
)

// Post ...
func (exe Execute) Post() (*Response, error) {
	bs, err := json.Marshal(exe)
	if err != nil {
		log.Infof("error marshalling the request: %v", err)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, compileURI, bytes.NewReader(bs))
	if err != nil {
		log.Errorf("error creating request: %v", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("error sending request from the default client: %v", err)
		return nil, err
	}

	bs, _ = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var codeResp Response
	if err = json.Unmarshal(bs, &codeResp); err != nil {
		log.Errorf("error unmarshalling: %v", err)
		return nil, err
	}

	if resp.StatusCode == 429 {
		log.Warnln("ran out of credits!")
		return &codeResp, ErrOutOfCredits
	} else if resp.StatusCode != 200 {
		log.Warnf("got a status code of %v and was expecing 200", resp.StatusCode)
		return &codeResp, ErrGotBadStatus

	}
	log.Infoln("successfully made post request to jquery got response", resp.StatusCode)
	return &codeResp, nil
}

// Handler is the rest api function handler for golang
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodOptions {
		// this is for cors
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var exe Execute
	err = json.Unmarshal(bs, &exe)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	exe.ClientID = os.Getenv("JDOODLE_ID")
	exe.ClientSecret = os.Getenv("JDOODLE_SECRET")

	resp, err := exe.Post()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bs, err = json.Marshal(*resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(bs)
}
