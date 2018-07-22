package bgaws

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/satori/go.uuid"
)

// Possible Question Difficulties
// beginner, easy, medium, hard, pro

var mQuestions []string

// GetQuestions gets a new question from Dynamo
func GetQuestions(difficulty string, amount int) (map[int]Question, error) {
	var qs []Question
	m := map[int]Question{}
	switch difficulty {
	//todo: add other difficulties
	case "medium":
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("us-east-1")},
		)

		if err != nil {
			return m, err
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

	if len(qs) == 0 {
		return m, errors.New("no questions were found")
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
		"72c0365e-c3cd-4cb7-8b87-4b6a018f2ecf",
		"7adee4ca-216b-497a-bcd9-c36d1676b211",
		"d6a8b122-5348-4183-9472-bfb28c8b2f42",
		"dcf864e3-28cb-420b-9f0f-63709c1e4ae8",
		"edd0e6f6-9f7d-43f1-a30f-3692ed926ffb",
		"f57b4c9b-3e0a-4685-85e6-ee909521a8dc",
		"f8f5f67d-88f1-4ef1-8e2f-bb5b93d1dba2",
	}
}
