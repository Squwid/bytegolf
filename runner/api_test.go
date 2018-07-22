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
