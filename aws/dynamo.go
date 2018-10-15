package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	usersTable = "bytegolf-users"
)

// GetUser todo
func GetUser(username, region, table string) (*User, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(region)},
	}))

	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})
	if err != nil {
		return &User{}, err
	}
	var u User
	err = dynamodbattribute.UnmarshalMap(result.Item, &u)
	if err != nil {
		return &User{}, err
	}
	return &u, nil
}

// CreateUser todo:
func CreateUser(user *User, region string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(region)},
	}))
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(*user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(usersTable),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

// UserExist checks to see if a user exists
func UserExist(username, region string) bool {
	user, _ := GetUser(username, region, usersTable)
	if len(user.Username) > 0 {
		return true
	}
	return false
}
