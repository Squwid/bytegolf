package question

import (
	"errors"

	"github.com/Squwid/bytegolf/firestore"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const collection = "questions"

// Question is a byte golf hole that gets store
type Question struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Question   string `json:"question"`
	TestCases  []Test `json:"tests"`
	TestCount  int    `json:"test_count"`
	Difficulty string `json:"difficulty"`
	Source     string `json:"source"`
	Live       bool   `json:"live"`
}

// Test is the tests that each compile request needs to pass in order to be correct
type Test struct {
	Input  string `json:"input"`
	Answer string `json:"answer"`
}

// Light is the same as a question but without any of the test cases
// so users cannot see the test questions, this is something that needs to change in the future
type Light struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Question   string `json:"question"`
	Live       bool   `json:"live"`
	Difficulty string `json:"difficulty"`
}

// ErrNil gets returned if a question is nil
var ErrNil = errors.New("given <nil> pointer")

// create creates a question, does not update or check for anything
// it will create a uuid for the id
func (q *Question) create() error {
	if q == nil {
		return ErrNil
	}
	if q.ID == "" {
		q.ID = uuid.New().String()
	}

	// if there are no test cases make it an empty list
	if q.TestCases == nil {
		q.TestCases = []Test{}
	}
	q.TestCount = len(q.TestCases)

	log.Infof("Creating new question %v (%v)", q.ID, q.Name)
	return firestore.StoreData(collection, q.ID, q)
}
