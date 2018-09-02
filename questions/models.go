package questions

import (
	"errors"
	"math/rand"
	"time"
)

// Error Variables
var (
	ErrNotEnoughQuestions = errors.New("not enough questions of that diffculty")
)

// Question is the type that is a question from the JSON and AWS API
type Question struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Question   string `json:"question"`
	Answer     string `json:"answer"`
	Difficulty string `json:"difficulty"`
	Source     string `json:"source"`
}

func randomize(questions []Question, amount int) map[int]Question {
	rand.Seed(time.Now().UTC().UnixNano())
	var counter int
	qs := make(map[int]Question)
	var used = []int{}
	for counter < amount {
		r := random(0, len(questions))
		if !contains(used, r) {
			counter++
			qs[counter] = questions[r]
			used = append(used, r)
		}
	}
	return qs
}

func contains(list []int, i int) bool {
	for _, item := range list {
		if item == i {
			return true
		}
	}
	return false
}

func random(min, max int) int {
	var r int
	r = min + rand.Intn(max)
	return int(r)
}
