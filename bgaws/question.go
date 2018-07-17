package bgaws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/satori/go.uuid"
)

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
