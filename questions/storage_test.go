package questions

import (
	"fmt"
	"testing"
)

func TestStoreQuestion(t *testing.T) {
	q := NewQuestion("testq", "TESTING", "answer", "medium", "att.net", "testq")
	err := q.Store(true)
	if err != nil {
		t.Logf("error storing question locally: %v\n", err)
		t.Fail()
	}
}

func TestGetLocalQuestions(t *testing.T) {
	qs := GetLocalQuestions()
	if len(qs) == 0 {
		t.Logf("expected more than 0 but got %v\n", len(qs))
		t.Fail()
	}
	fmt.Println("QS:", qs)
}
