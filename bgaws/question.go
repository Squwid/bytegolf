package bgaws

import (
	"errors"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/satori/go.uuid"
)

var easyQs = []string{"5f8d6ab3-bf6b-4347-bcfa-305a5ec4cb7e", "3da0bd40-799b-4357-b694-86e1ecb93e4e"}
var mediumQs = []string{"0f9ad9f1-1bda-487b-be07-fe691d1a056b", "309c6a85-d18d-4ba2-8b0b-e928107597ae", "3da0bd40-799b-4357-b694-86e1ecb93e4e", "72c0365e-c3cd-4cb7-8b87-4b6a018f2ecf", "7adee4ca-216b-497a-bcd9-c36d1676b211", "dcf864e3-28cb-420b-9f0f-63709c1e4ae8", "dcf864e3-28cb-420b-9f0f-63709c1e4ae8", "f8f5f67d-88f1-4ef1-8e2f-bb5b93d1dba2"}

// Error Variables
var (
	ErrNotEnoughQuestions = errors.New("not enough questions exist for this difficulty")
)

// GetQuestions gets a new question from Dynamo
func GetQuestions(difficulty string, amount int) (map[int]Question, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	var qs []Question
	m := map[int]Question{}
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return m, err
	}
	switch difficulty {
	case "easy":
		if amount > len(easyQs) {
			return m, ErrNotEnoughQuestions
		}
		// var previous = []int{}
		// for i := 1; i <= amount; i++ {
		// 	var r
		// 	m[i] =
		// }

	case "medium":
		if amount > len(mediumQs) {
			return m, ErrNotEnoughQuestions
		}

		qIDs := createMQs(amount)
		if len(qIDs) < amount {
			return m, errors.New("not enough questions exist for this difficulty")
		}

		svc := dynamodb.New(sess)
		for _, id := range qIDs {
			result, err := svc.GetItem(&dynamodb.GetItemInput{
				TableName: aws.String(qsTable),
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
	}

	for i, q := range qs {
		m[i+1] = q
	}
	return m, nil
}

// Store stores a question inside of AWS DynamoDB
func (q *Question) Store() error {
	qID, _ := uuid.NewV4()
	q.ID = qID.String()
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String("us-east-1")},
	}))

	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(*q)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(qsTable),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

// ["72c0365e-c3cd-4cb7-8b87-4b6a018f2ecf","7adee4ca-216b-497a-bcd9-c36d1676b211","d6a8b122-5348-4183-9472-bfb28c8b2f42"]
func createMQs(amount int) []string {
	return []string{
		"3da0bd40-799b-4357-b694-86e1ecb93e4e",
		"5f8d6ab3-bf6b-4347-bcfa-305a5ec4cb7e",
		"72efd1af-9e88-423e-a23f-e0f0612eea5e",
		"72c0365e-c3cd-4cb7-8b87-4b6a018f2ecf",
		"7adee4ca-216b-497a-bcd9-c36d1676b211",
		"d6a8b122-5348-4183-9472-bfb28c8b2f42",
		"dcf864e3-28cb-420b-9f0f-63709c1e4ae8",
		"edd0e6f6-9f7d-43f1-a30f-3692ed926ffb",
		"f57b4c9b-3e0a-4685-85e6-ee909521a8dc",
		"f8f5f67d-88f1-4ef1-8e2f-bb5b93d1dba2",
	}
}

func random(min, max int) int {
	var r int
	r = min + rand.Intn(max)
	return int(r)
}
