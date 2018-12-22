package questions

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/satori/go.uuid"
)

// Store stores a question locally if the local bool is true
func (q *Question) Store(local bool) error {
	if local {
		return q.storeLocal()
	}
	return errors.New("Not storing question because local is false")
}

func (q *Question) storeLocal() error {
	var path = fmt.Sprintf("./questions/questions/")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	f, err := os.Create(fmt.Sprintf("%s%s.txt", path, q.Link))
	if err != nil {
		return err
	}
	defer f.Close()

	store, err := json.Marshal(q)
	if err != nil {
		return err
	}

	_, err = f.Write(store)
	if err != nil {
		return err
	}

	return nil
}

// Store stores a question inside of AWS DynamoDB
func (q *Question) storeAWS() error {
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
