package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"golang.org/x/crypto/bcrypt"
)

// User is the user for bytegolf, stored in an aws table
type User struct {
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
	Role        string `json:"role"`

	Created   string `json:"created"`
	LastLogin string `json:"lastLogin"`
}

// NewUser creates a new user using an email as the primary key. The password is encrypted and panics on error
func NewUser(email, displayName, password string) *User {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err) // panics the error because this will only happen if mem is full
	}
	return &User{
		Email:       email,
		DisplayName: displayName,
		Password:    string(encrypted),
		Role:        RoleUser,
		Created:     time.Now().String(),
		LastLogin:   time.Now().String(),
	}
}

// GetUser retrieves a user using a specific string from where they are stored
func GetUser(email string) (*User, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(DefaultRegion)},
	})
	if err != nil {
		return nil, err
	}
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(usersTable),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	var u User
	err = dynamodbattribute.UnmarshalMap(result.Item, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Store stores a user into the database
func (user *User) Store() error {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(DefaultRegion)},
	})
	if err != nil {
		return err
	}
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
	return err
}

// UserExist checks to see if a user exists. It returns the user as well for use in caching purposes
func UserExist(email string) (bool, *User) {
	user, err := GetUser(email)
	if err != nil {
		panic(err) // panics because this would happen only if the aws is down
	}
	if email == user.Email {
		return false, user
	}
	return true, user
}
