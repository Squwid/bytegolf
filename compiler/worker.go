package main

import (
	"fmt"
	"time"

	"github.com/Squwid/bytegolf/lib/api"
	"github.com/sirupsen/logrus"
)

const workerCount = 5
const workerBacklog = workerCount * 4

var workerQueue = make(chan WorkerInput, workerBacklog)

type WorkerInput struct {
	Language   *api.LanguageDB
	Test       *api.TestDB
	Submission *api.SubmissionDB
}

type Worker struct {
	ID string
}

func init() {
	for i := 0; i < workerCount; i++ {
		worker := &Worker{}
		go worker.Start()
	}
}

func (worker *Worker) Start() {
	worker.ID = api.RandomString()
	logger := logrus.WithField("ID", worker.ID)
	logger.Info("Worker started")
	defer logger.Warnf("Worker ded")

	for {
		input := <-workerQueue

		time.Sleep(1 * time.Second)
		fmt.Printf("Worker %s is working on %s\n", worker.ID, input.Submission.ID)
	}
}
