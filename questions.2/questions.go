package question

import (
	"context"
	"errors"

	"github.com/Squwid/bytegolf/firestore"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

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

type Test struct {
	Input  string `json:"input"`
	Answer string `json:"answer"`
}

// NewQuestion returns a new Question after generating a uuid
func NewQuestion() Question {
	return Question{ID: uuid.New().String()}
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
	return firestore.StoreData("questions", q.ID, q)
}

// GetLiveQuestions gets a list of questions that have the Live bool
func GetLiveQuestions() ([]Question, error) {
	var qs = []Question{}
	ctx := context.Background()
	iter := firestore.Client.Collection("questions").Where("Live", "==", true).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var q Question
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
