package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Squwid/bytegolf/runner"

	"github.com/Squwid/bytegolf/aws"
)

func TestScoring(t *testing.T) {
	testQ := aws.Question{
		ID:         "2",
		Name:       "test question",
		Question:   "print `hello world!`",
		Answer:     "hello world!",
		Difficulty: "easy",
	}

	codeBodyGo := `
	package main
	import "fmt"



	func main() {
	fmt.Println("hello world!")

	tkjakdsjf
	}
	`
	c := runner.NewClient()
	sub := runner.NewCodeSubmission("bwhitelaw24", "1", "main.go", runner.LangGo, codeBodyGo, c, nil)
	resp := runner.CodeResponse{
		UUID:       "1",
		Output:     "hello world!",
		StatusCode: 200,
		Memory:     "some mem value",
		CPUTime:    "NONE",
	}
	correct := checkResponse(&resp, &testQ)
	if !correct {
		fmt.Println("QUESTION WAS INCORRECT")
		t.Fail()
	}

	score := Score(sub, &testQ)
	t.Log("GOT", score, "FROM QUESTION", testQ.Name)

	{
		return
		// COUNT DEBUGGER
		var c uint
		for _, l := range sub.Script {
			if len(strings.TrimSpace(string(l))) == 0 {
				continue
			} else {
				t.Log("*** CHAR:", string(l))
			}
			c++
		}
	}
}
