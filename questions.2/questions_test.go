package question

import (
	"testing"
)

func TestQuestions(t *testing.T) {
	var q = &Question{
		Name:     "ben",
		Question: "what is 2+2",
		TestCases: []Test{
			{
				Input:  "123",
				Answer: "abc",
			},
		},
		TestCount:  1,
		Difficulty: "easy",
		Source:     "ben.com",
		Live:       true,
	}
	err := q.create()
	if err != nil {
		t.Fatalf("%v\n", err)
	}
}

func TestRetreiveQuestions(t *testing.T) {
	// var qs = []Question{}
	q, err := GetLiveQuestions()
	if err != nil {
		t.Fatalf("Error getting live questions: %v", err)
	}

	t.Logf("%v", q)
}
