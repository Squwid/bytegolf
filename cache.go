package main

import (
	"github.com/Squwid/bytegolf/aws"
	"golang.org/x/crypto/bcrypt"
)

// usersCache holds the 'maxCache' most recently accessed users incase of wrong passwords and such
// to deal with how many read/writes
var usersCache []*aws.User

const maxCache = 10

func cacheUsers(users ...*aws.User) {
	for _, u := range users {
		index, exist := userInCache(u.Email)
		if exist {
			// if the user exists change the order of the cache
			usersCache = append(usersCache[:index], usersCache[index+1:]...)
			usersCache = append(usersCache, u)
			return
		}
		if len(usersCache) >= maxCache {
			usersCache = append(usersCache[:0], usersCache[1:]...)
		}
		usersCache = append(usersCache, u)
	}
}

// userInCache checks to see if the user is in the users cache, and if they are the index of that element gets returned
func userInCache(email string) (int, bool) {
	for i, u := range usersCache {
		if email == u.Email {
			return i, true
		}
	}
	return 0, false
}

// getAwsUser does the same thing that aws.GetUser does but it checks the cache first to save as many readwrites from aws as possible
func getAwsUser(email string) (*aws.User, error) {
	index, exist := userInCache(email)
	if exist {
		return usersCache[index], nil
	}
	return aws.GetUser(email)
}

// tryLogin tries an email and password and checks to see if its correct. It uses user caching
// incase the user tries multiple logins. Returns errors if aws does not act as intended
func tryLogin(email, password string) (bool, error) {
	user, err := getAwsUser(email)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logger.Printf("%s tried to login incorrectly\n", email)
		return false, nil
	}
	return true, nil
}
