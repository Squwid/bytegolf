package compiler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

const compileURI = ""

type Execute struct {
	Script       string `json:"script"`
	Language     string `json:"language"`
	VersionIndex string `json:"versionIndex"`

	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

// Post ...
func (exe Execute) Post() error {
	exe.ClientID = os.Getenv("BG_CLIENT_ID")
	exe.ClientSecret = os.Getenv("BG_CLIENT_SECRET")

	bs, err := json.Marshal(exe)
	if err != nil {
		log.Infof("error marshalling the request: %v", err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, compileURI, bytes.NewReader(bs))
	if err != nil {
		log.Errorf("error creating request: %v", err)
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	log.Info("success")
	log.Println(resp)
	return nil
}
