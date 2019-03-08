package questions

import (
	"testing"
)

func TestStoreQuestion(t *testing.T) {
	q := NewQuestion("testq", "TESTING", "answer", "medium", "att.net", "testq")
	err := q.Store()
	if err != nil {
		t.Errorf("error storing question locally: %v\n", err)
	}
}

func TestGetAllQuestions(t *testing.T) {
	q := NewQuestion("testq", "TESTING", "answer", "medium", "att.net", "testq")
	err := q.Store()
	if err != nil {
		t.Fatalf("error storing question locally: %v\n", err)
	}

	qs, err := GetLocalQuestions()
	if err != nil {
		t.Errorf("error getting all questions : %v\n", err)
	}
	if len(qs) < 1 {
		t.Errorf("expected more than 1 question but got %v\n", len(qs))
	}
}

func TestLiveQuestions(t *testing.T) {
	q := NewQuestion("testq", "TESTING", "answer", "medium", "att.net", "testq")
	err := q.Store()
	if err != nil {
		t.Fatalf("error storing question locally: %v\n", err)
	}

	err = q.Deploy(1)
	if err != nil {
		t.Fatalf("error deploying to hole 1 : %v\n", err)
	}
	lives, err := GetLiveQuestions()
	if err != nil {
		t.Errorf("error getting live questions : %v\n", err)
	}
	if len(lives) < 1 {
		t.Errorf("expected more than 1 live questions\n")
	}
}
