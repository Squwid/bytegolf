package runner

import (
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
		t.Log("c.ID:", c.ID)
		t.Log("error creating new client")
		t.Fail()
	}
	sub := NewCodeSubmission("bwhitelaw24", "343", "main.go", LangGo, codeBody, c)
	resp, err := sub.Send()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(resp)
}

func TestStoreSubmissionLocal(t *testing.T) {
	c := NewClient()
	sub := NewCodeSubmission("bwhitela4", "343", "main.go", LangGo, codeBody, c)
	err := sub.storeLocal()
	if err != nil {
		t.Logf("Error storing local: %v\n", err)
		t.Fail()
	}
}

func TestStoreResponseLocal(t *testing.T) {
	c := NewClient()
	sub := NewCodeSubmission("bwhitela4", "343", "main.go", LangGo, codeBody, c)
	resp, err := sub.Send()
	if err != nil {
		t.Logf("error sending submission: %v\n", err)
		t.Fail()
	}

	err = resp.storeLocal()
	if err != nil {
		t.Logf("error storing response locally: %v\n", err)
		t.Fail()
	}
}
