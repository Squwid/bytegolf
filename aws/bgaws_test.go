package aws

import (
	"testing"
)

func TestDynamoGet(t *testing.T) {
	_, err := GetUser("phil")
	if err != nil {
		t.Log("Could not get user ben", err)
		t.Fail()
	}
}

func TestDynamoPost(t *testing.T) {
	ben := User{
		Username: "username",
		Password: "password",
		Role:     "user",
	}

	err := CreateUser(&ben)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}
