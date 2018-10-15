package aws

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/satori/go.uuid"
)

// Error Variables
var (
	ErrNotEnoughQuestions = errors.New("not enough questions of that diffculty")
)

// Question is the type that gets stored as a question in dynamodb
type Question struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Question   string `json:"question"`
	Answer     string `json:"answer"`
	Difficulty string `json:"difficulty"`
	Source     string `json:"source"`
}

// FIXME: Hard coded questions need to be generated from either an API or library or local (local should be an option)

var easyQs = []string{"5f8d6ab3-bf6b-4347-bcfa-305a5ec4cb7e", "3da0bd40-799b-4357-b694-86e1ecb93e4e"}
var mediumQs = []string{"0f9ad9f1-1bda-487b-be07-fe691d1a056b", "309c6a85-d18d-4ba2-8b0b-e928107597ae", "3da0bd40-799b-4357-b694-86e1ecb93e4e", "72c0365e-c3cd-4cb7-8b87-4b6a018f2ecf", "7adee4ca-216b-497a-bcd9-c36d1676b211", "dcf864e3-28cb-420b-9f0f-63709c1e4ae8", "dcf864e3-28cb-420b-9f0f-63709c1e4ae8", "f8f5f67d-88f1-4ef1-8e2f-bb5b93d1dba2"}

// Store stores a question inside of AWS DynamoDB
func (q *Question) Store() error {
	qID, _ := uuid.NewV4()
	q.ID = qID.String()
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(DefaultRegion)},
	}))

	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(*q)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(questionsTable),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

// ["72c0365e-c3cd-4cb7-8b87-4b6a018f2ecf","7adee4ca-216b-497a-bcd9-c36d1676b211","d6a8b122-5348-4183-9472-bfb28c8b2f42"]
func createMQs(amount int) []string {
	switch amount {
	case 1:
		return []string{"3da0bd40-799b-4357-b694-86e1ecb93e4e"}
	case 3:
		return []string{
			"3da0bd40-799b-4357-b694-86e1ecb93e4e",
			"0f9ad9f1-1bda-487b-be07-fe691d1a056b",
			"72efd1af-9e88-423e-a23f-e0f0612eea5e",
		}
	}
	return []string{
		"3da0bd40-799b-4357-b694-86e1ecb93e4e",
		"5f8d6ab3-bf6b-4347-bcfa-305a5ec4cb7e",
		"72efd1af-9e88-423e-a23f-e0f0612eea5e",
		"72c0365e-c3cd-4cb7-8b87-4b6a018f2ecf",
		"7adee4ca-216b-497a-bcd9-c36d1676b211",
		"309c6a85-d18d-4ba2-8b0b-e928107597ae",
		"dcf864e3-28cb-420b-9f0f-63709c1e4ae8",
		"edd0e6f6-9f7d-43f1-a30f-3692ed926ffb",
		"f8f5f67d-88f1-4ef1-8e2f-bb5b93d1dba2",
	}
}

// GetQuestionsLocal gets questions stored locally inside of this folder, incase of no AWS
func GetQuestionsLocal(amount int, difficulty string) (map[int]Question, error) {
	tempQs := []Question{}
	questions := []Question{}
	file, err := ioutil.ReadFile("questions.json")
	if err != nil {
		fmt.Println("error at 1 :", err)
		return map[int]Question{}, err
	}
	err = json.Unmarshal(file, &questions)
	if err != nil {
		fmt.Println("error at 2 :", err)
		return map[int]Question{}, err
	}
	for _, q := range questions {
		if q.Difficulty == difficulty {

			tempQs = append(tempQs, q)
		}
	}
	if len(tempQs) < amount {
		return map[int]Question{}, fmt.Errorf("not enough %s questions, wanted %v got %v", difficulty, amount, len(tempQs))
	}
	return randomize(tempQs, amount), nil
}

// GetQuestionsDynamo gets questions from a table inside of DynamoDB
func GetQuestionsDynamo(amount int, difficulty string) (map[int]Question, error) {
	var qs []Question
	m := make(map[int]Question)
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(DefaultRegion)}),
	)

	switch difficulty {
	case "easy":
		return m, ErrNotEnoughQuestions

	case "medium":
		if amount > 9 {
			return m, ErrNotEnoughQuestions
		}

		qIDs := createMQs(amount)
		svc := dynamodb.New(sess)
		for _, id := range qIDs {
			result, err := svc.GetItem(&dynamodb.GetItemInput{
				TableName: aws.String(questionsTable),
				Key: map[string]*dynamodb.AttributeValue{
					"id": {
						S: aws.String(id),
					},
				},
			})
			if err != nil {
				return m, err
			}
			var q Question
			err = dynamodbattribute.UnmarshalMap(result.Item, &q)
			if err != nil {
				return m, err
			}
			qs = append(qs, q)
		}

	case "hard":
		return m, ErrNotEnoughQuestions
	}

	return randomize(qs, amount), nil
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
