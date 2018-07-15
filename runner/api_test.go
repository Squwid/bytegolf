package runner

import (
	"testing"
)

func TestSubmit(t *testing.T) {
	c := NewClient()
	if len(c.ID) == 0 {
		t.Log("c.ID:", c.ID)
		t.Log("error creating new client")
		t.Fail()
	}
	code := `package main

	import "fmt"

	func main(){
		fmt.Println("hello world")
	}
	`
	sub := NewCodeSubmission(LangGo, code, c)
	resp, err := sub.Send()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(resp)
}
