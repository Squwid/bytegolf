package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const uriAPI = "https://api.jdoodle.com/execute"

// Send todo
func (s *CodeSubmission) Send() (*CodeResponse, error) {
	var codeResponse CodeResponse

	reqBody, err := json.Marshal(*s)
	if err != nil {
		return &CodeResponse{}, err
	}

	fmt.Println("Body:", string(reqBody))

	req, err := http.NewRequest(http.MethodPost, uriAPI, bytes.NewBuffer(reqBody))
	if err != nil {
		return &CodeResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &CodeResponse{}, err
	}

	err = json.Unmarshal(body, &codeResponse)
	if err != nil {
		return &CodeResponse{}, err
	}

	return &codeResponse, nil
}
