package questions

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Squwid/bytegolf/runner"
)

func TestScoring(t *testing.T) {
	testQ := NewQuestion("test question", "print `hello world!`", "hello world!", "easy", "source", "hw")

	codeBodyGo := `
	package main
	import "fmt"



	func main() {
	fmt.Println("hello world!")

	tkjakdsjf
	}
	`
	c := runner.NewClient()
	sub := runner.NewCodeSubmission("bwhitelaw24", "1", "main.go", runner.LangGo, codeBodyGo, c)
	resp := runner.CodeResponse{
		UUID:       "1",
		Output:     "hello world!",
		StatusCode: 200,
		Memory:     "some mem value",
		CPUTime:    "NONE",
	}
	correct := testQ.Check(&resp)
	if !correct {
		fmt.Println("QUESTION WAS INCORRECT")
		t.Fail()
	}

	score := testQ.Score(sub)
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
