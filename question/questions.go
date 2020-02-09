package question

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	fs "github.com/Squwid/bytegolf/firestore"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
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
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
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
	return fs.StoreData(collection, q.ID, q)
}

// GetQuestion gets a single question using an ID, returns a nil question if the question is not found
func GetQuestion(id string) (*Question, error) {
	var qs = []Question{}
	ctx := context.Background()
	iter := fs.Client.Collection(collection).Where("ID", "==", id).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Errorf("Error getting question with id %v: %v", id, err)
			return nil, err
		}

		var q Question
		err = mapstructure.Decode(doc.Data(), &q)
		if err != nil {
			log.Errorf("Error decoding question with id %v: %v ", id, err)
			return nil, err
		}
		qs = append(qs, q)
	}
	if len(qs) > 1 {
		log.Warnf("Request to list hole %s returned %s questions", id, len(qs))
	}

	if len(qs) == 0 {
		return nil, nil
	}

	return &qs[0], nil
}

// TransformToLight takes a question and transforms it to a light question
// which hides things like test cases for the user
func (q Question) TransformToLight() Light {
	return Light{
		ID:         q.ID,
		Name:       q.Name,
		Question:   q.Question,
		Live:       q.Live,
		Difficulty: q.Difficulty,
	}
}

// listQuestions gets a list of questions that have the Live bool
func listQuestions(onlyLive bool) ([]Light, error) {
	var qs = []Light{}
	ctx := context.Background()

	// if onlyLive return only live questions, if not return every question
	var iter *firestore.DocumentIterator
	if onlyLive {
		// only live questions
		iter = fs.Client.Collection(collection).Where("Live", "==", true).Documents(ctx)
	} else {
		// all questions
		iter = fs.Client.Collection(collection).Documents(ctx)
	}

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var q Light
		err = mapstructure.Decode(doc.Data(), &q)
		if err != nil {
			log.Errorf("error decoding object: %v", err)
		} else {
			log.Debugf("got data back, parsing: %s", doc.Data())
			qs = append(qs, q)
		}
	}
	return qs, nil
}
