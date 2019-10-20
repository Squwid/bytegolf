package compiler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Squwid/bytegolf/firestore"
	"github.com/Squwid/bytegolf/secrets"
	"github.com/Squwid/bytegolf/sess"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const compileURI = "https://api.jdoodle.com/v1/execute"

var jdoodleClient *secrets.Client

func init() {
	jdoodleClient = secrets.Must(secrets.GetClient("JDOODLE")).(*secrets.Client)
}

// Execute ...
type Execute struct {
	Script       string `json:"script"`
	Language     string `json:"language"`
	VersionIndex string `json:"versionIndex"`
	HoleID       string `json:"holeId"`

	// TODO: remove these from here and put them somewhere else
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

// Response ...
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
func (exe Execute) Post(s *sess.Session) (*Response, error) {
	// loggedIn, err := sess.LoggedIn(req)

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

	// should only store if the request is 200 for now, but could come back later and move this somewhere else
	// as a q if credits run out or some other error
	go func(exe Execute, r Response, bgid string) {
		var c = TotalStore{
			Exe:     exe,
			Resp:    codeResp,
			BGID:    bgid,
			Correct: true,
			HoleID:  exe.HoleID,
		}
		uid := uuid.New().String()
		err := firestore.StoreData("executes", uid, c)
		if err != nil {
			log.Errorf("error storing fire data store (%s): %v", uid, err)
			return
		}
		log.Infof("Saved firestore data for %s at %s", bgid, uid)
	}(exe, codeResp, s.BGID)
	log.Infoln("successfully made post request to jquery got response", resp.StatusCode)
	return &codeResp, nil
}

// TotalStore ...
type TotalStore struct {
	Exe     Execute  `json:"submission"`
	Resp    Response `json:"response"`
	BGID    string   `json:"bgid"`
	Correct bool     `json:"correct"`
	HoleID  string   `json:"holeId"`
}

// Handler is the rest api function handler for golang
func Handler(w http.ResponseWriter, r *http.Request) {
	loggedIn, s, err := sess.LoggedIn(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("error checking to see if a user is signed in: %v", err)
		return
	}
	if !loggedIn {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprintf(`{"error": "unauthorized"}`)))
		return
	}
	if s == nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("error: session was blank")
		return
	}
	// the user is logged in

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS,GET")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodOptions {
		// this is for cors
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method == http.MethodGet {
		getExecutes(w, r, s)
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
	// change to secrets manager from environmental variables
	exe.ClientID = jdoodleClient.Client
	exe.ClientSecret = jdoodleClient.Secret

	resp, err := exe.Post(s)
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
