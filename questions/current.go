package questions

import (
	"log"
	"sync"
	"time"
)

// question variables, such as the lock for reading and writing
var (
	qs     = map[int]*Question{}
	qMutex = &sync.Mutex{}
)

func init() {
	qs = Must(MapLiveQuestions())
	// go updateQuestions(24 * time.Hour)
}

// GetQuestion gets a specific question from the map of questions, if it does not exist, an error is returned
func GetQuestion(holeID string) Question {
	qMutex.Lock()
	defer qMutex.Unlock()

	for _, q := range qs {
		if q.ID == holeID {
			return *q
		}
	}
	return Question{}
}

// updateQuestions updates the questions every certain amount of time
func updateQuestions(often time.Duration) {
	for range time.Tick(often) {
		UpdateQuestions()
	}
}

// UpdateQuestions updates the users questions from the files. This is so they are stored in a map rather than local files
func UpdateQuestions() {
	s, err := MapLiveQuestions()
	if err != nil {
		log.Println("error updating qs", err)
	} else {
		qMutex.Lock()
		qs = s
		qMutex.Unlock()
	}
}

// Must takes in a map and a error, and if an error has occurred it will panic
func Must(qs map[int]*Question, err error) map[int]*Question {
	if err != nil {
		panic(err)
	}
	return qs
}

// GetAllQuestions returns the map of questions
func GetAllQuestions() map[int]*Question {
	newq := map[int]*Question{}
	qMutex.Lock()
	for k, v := range qs {
		newq[k] = v
	}
	qMutex.Unlock()
	return newq
}
