package api

import (
	"math/rand"
	"sync"
	"time"
)

var randMutex = &sync.Mutex{}

const charset = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`

// RandomString generates a 6 character/number string.
func RandomString(length int) string {
	randMutex.Lock()
	defer randMutex.Unlock()

	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
