package runner

import (
	"encoding/json"
	"net/http"
	"testing"
)

const codeBody = `package main
import "fmt"

func main() {
	fmt.Println("hello world!")
}
	`

func TestSubmit(t *testing.T) {
	c := NewClient()
	if len(c.ID) == 0 {
		t.Fatalf("error creating new client w/id: %s\n", c.ID)
	}
	sub := NewCodeSubmission("bwhitelaw24", "343", "", "main.go", LangGo, codeBody, c, nil)
	resp, err := sub.Send(false)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("got bad status : %v", resp.StatusCode)
	} else {
		t.Logf(resp.Output)
		bs, _ := json.Marshal(*resp)
		t.Log(string(bs))
	}
}

/*
func TestStoreSubmissionLocal(t *testing.T) {
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
}
/*
func TestStoreResponseLocal(t *testing.T) {
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
		t.Errorf("error storing response locally: %v\n", err)
	}
	if _, err := os.Stat("./resp/bwhitelaw/1"); os.IsNotExist(err) {
		t.Errorf("file was not created")
	}
}
*/
