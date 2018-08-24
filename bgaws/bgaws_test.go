package bgaws

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDynamoGet(t *testing.T) {
	_, err := GetUser("phil")
	if err != nil {
		t.Log("Could not get user ben", err)
		t.Fail()
	}
}

func TestDynamoPost(t *testing.T) {
	ben := User{
		Username: "username",
		Password: "password",
		Role:     "user",
	}

	err := CreateUser(&ben)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

// so far we have 1-5 project euler problems inserted into the database
func TestCreateQuestion(t *testing.T) {
	question := Question{
		Name: "1000-digit Fibonacci number",
		Question: `The Fibonacci sequence is defined by the recurrence relation:

		Fn = Fn−1 + Fn−2, where F1 = 1 and F2 = 1.
		Hence the first 12 terms will be:
		
		F1 = 1
		F2 = 1
		F3 = 2
		F4 = 3
		F5 = 5
		F6 = 8
		F7 = 13
		F8 = 21
		F9 = 34
		F10 = 55
		F11 = 89
		F12 = 144
		The 12th term, F12, is the first term to contain three digits.
		
		What is the index of the first term in the Fibonacci sequence to contain 1000 digits?`,
		Answer:     "4782",
		Difficulty: "hard",
		Source:     "ben",
	}
	err := question.Store()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestPrintQuestion(t *testing.T) {
	q := Question{
		Name: "Summation of primes",
		Question: `The sum of the primes below 10 is 2 + 3 + 5 + 7 = 17.

		Find the sum of all the primes below two million.`,
		Answer:     "104743",
		Difficulty: "medium",
		Source:     "projecteuler.net",
	}
	bs, _ := json.Marshal(q)
	fmt.Println(string(bs))
}

func TestGetQuestion(t *testing.T) {
	qs, err := GetQuestions("medium", 3)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Logf("found %v from question medium\n", len(qs))
	t.Logf("\t%v\n", qs)
}
