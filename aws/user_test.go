package aws

import (
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	var newUserEmail = "b@t.net"
	var newUserName = "JJJ"
	var newUserFirst = "Jypsy"
	u := NewUser(newUserEmail, newUserName, newUserFirst)
	if err := u.Store(); err != nil {
		t.Logf("error occurred storing user: %v\n", err)
		t.Fail()
	}
	time.Sleep(3 * time.Second)
	user, err := GetUser(u.Email)
	if err != nil {
		t.Logf("error occurred getting user: %v\n", err)
		t.Fail()
	}
	if user.Email != u.Email {
		t.Logf("users email was %s when it should be %s\n", user.Email, u.Email)
		t.Fail()
	}
}
