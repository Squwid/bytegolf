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
	exe.ClientID = os.Getenv("BG_CLIENT_ID")
	exe.ClientSecret = os.Getenv("BG_CLIENT_SECRET")

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
