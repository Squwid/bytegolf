package question

import "testing"

func TestAddQuestion(t *testing.T) {
	q := NewQuestion()
	q.Name = "Bens First Question"
	q.Question = "Whats three plus 3"
	q.Difficulty = "easy"
	q.Source = "https://google.com"
	q.Live = true
	err := q.create()
	if err != nil {
		t.Errorf("error storing question: %v", err)
		return
	}
}
