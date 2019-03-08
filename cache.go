package main

import (
	"sync"

	"github.com/Squwid/bytegolf/aws"
)

// usersCache holds the 'maxCache' most recently accessed users incase of wrong passwords and such
// to deal with how many read/writes
var (
	usersCache []aws.User
	usersMutex = &sync.Mutex{}
)

const maxCache = 10

func cacheUsers(users ...*aws.User) {
	for _, u := range users {
		index, exist := userInCache(u.Email)
		if exist {
			// if the user exists change the order of the cache
			usersCache = append(usersCache[:index], usersCache[index+1:]...)
			usersCache = append(usersCache, *u)
			return
		}
		if len(usersCache) >= maxCache {
			usersCache = append(usersCache[:0], usersCache[1:]...)
		}
		usersCache = append(usersCache, *u)
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
func getAwsUser(email string) (aws.User, error) {
	// make the cache threadsafe
	usersMutex.Lock()
	defer usersMutex.Unlock()

	index, exist := userInCache(email)
	if exist {
		return usersCache[index], nil
	}
	logger.Printf("\t*** Getting aws user\n")
	user, err := aws.GetUser(email)
	cacheUsers(user)
	return *user, err
}
