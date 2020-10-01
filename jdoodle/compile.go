package jdoodle

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Squwid/bytegolf/models"
)

const jdoodleURL = "https://api.jdoodle.com/v1/execute"

var jdoodleClientID string
var jdoodleClientSecret string

// Errors
var (
	ErrOutOfCredits = errors.New("Out of credits error")
)

func init() {
	jdoodleClientID = os.Getenv("JDOODLE_CLIENT")
	jdoodleClientSecret = os.Getenv("JDOODLE_SECRET")
}

// SendJdoodle takes a CompileInput and sends it through the jdoodle online compiler
func SendJdoodle(input models.CompileInput) (*models.CompileOutput, error) {
	// Body with the clientid and secret
	body := struct {
		models.CompileInput
		ClientID     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
	}{
		input,
		jdoodleClientID,
		jdoodleClientSecret,
	}

	bs, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Send POST request to Jdoodle
	resp, err := http.DefaultClient.Post(jdoodleURL, "application/json", bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	// 429 - Out of credits
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrOutOfCredits
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error unexpected status code %v", resp.StatusCode)
	}

	// Got response and its 200, parse
	var out models.CompileOutput
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
