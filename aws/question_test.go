package aws

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestCreateQuestions(t *testing.T) {
	q1 := Question{
		ID:         "1",
		Name:       "Sample Question 1",
		Question:   "What is the first question ever?",
		Answer:     "That one",
		Difficulty: "hard",
	}
	q2 := Question{
		ID:         "2",
		Name:       "Sample Question 2",
		Question:   "What was the answer to the second question?",
		Answer:     "Whats the second question",
		Difficulty: "easy",
	}
	var qs = []Question{q1, q2}
	bs, err := json.Marshal(qs)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(string(bs))
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

func TestGetQuestions(t *testing.T) {
	qs, err := GetQuestionsDynamo(4, "medium")
	if err != nil {
		t.Logf("error getting questions from dynamo: %v\n", err)
		t.Fail()
	}
	if len(qs) != 4 {
		t.Logf("expecting 4 questions got %v\n", len(qs))
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
