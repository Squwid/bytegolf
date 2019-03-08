package runner

import (
	"os"
	"testing"
)

func TestGetPlayerSubmissions(t *testing.T) {
	sub := &CodeSubmission{
		UUID:         "1",
		Script:       codeBody,
		Language:     LangGo,
		VersionIndex: "0",
		ID:           "abc",
		Secret:       "def",
		Info: &FileInfo{
			Name: "Ben",
			User: "bwhitelaw",
			Hole: "3",
		},
	}

	err := sub.storeLocal()
	if err != nil {
		t.Errorf("Error storing local: %v\n", err)
	}
	if _, err := os.Stat("./subs/bwhitelaw/1"); os.IsNotExist(err) {
		t.Errorf("file was not created")
	}

	subs, err := GetPlayerSubmissions("bwhitelaw")
	if err != nil {
		t.Fatal(err)
	}
	if len(subs) < 1 {
		t.Errorf("expected at least 2 subs but got %v\n", err)
	}
}

func TestGetPlayerResponses(t *testing.T) {
	resp := &CodeResponse{
		UUID:       "1",
		Output:     "hello world",
		StatusCode: 200,
		Info: &FileInfo{
			Name: "Ben",
			User: "bwhitelaw",
			Hole: "3",
		},
	}

	err := resp.storeLocal()
	if err != nil {
		t.Errorf("Error storing local: %v\n", err)
	}
	if _, err := os.Stat("./subs/bwhitelaw/1"); os.IsNotExist(err) {
		t.Errorf("file was not created")
	}

	resps, err := GetPlayerResponses("bwhitelaw")
	if err != nil {
		t.Fatal(err)
	}
	if len(resps) < 1 {
		t.Errorf("expected at least 2 subs but got %v\n", err)
	}
}
