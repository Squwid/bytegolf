package bgaws

import "testing"

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

// so far we have 1-5 project euler problems inserted into the database
func TestCreateQuestion(t *testing.T) {
	question := Question{
		Name: "Smallest multiple",
		Question: `2520 is the smallest number that can be divided by each of the numbers from 1 to 10 without any remainder.

		What is the smallest positive number that is evenly divisible by all of the numbers from 1 to 20?`,
		Answer:     "232792560",
		Difficulty: "medium",
		Source:     "projecteuler.net",
	}
	err := question.Store()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}
