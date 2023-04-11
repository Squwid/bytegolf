package processor

import (
	"fmt"
	"sync"
	"time"

	"github.com/Squwid/bytegolf/lib/docker"
	"github.com/Squwid/bytegolf/lib/log"
	"github.com/pkg/errors"
)

const (
	timeout = 10 * time.Second

	jobBacklog  = 5000
	bytesToRead = 4096
	workerCount = 4
)

var JobQueue = make(chan *Job, jobBacklog)

type Worker struct {
	ID string

	lock *sync.Mutex
}

var WorkerPool [workerCount]*Worker

func Init() {
	for i := 0; i < workerCount; i++ {
		worker := NewWorker(i)
		WorkerPool[i] = worker
		go worker.Start()
	}
}

func NewWorker(id int) *Worker {
	return &Worker{
		ID:   fmt.Sprintf("w%v", id),
		lock: &sync.Mutex{},
	}
}

func (worker *Worker) Start() {
	workerLogger := log.GetLogger().WithField("Worker", worker.ID)
	workerLogger.Debugf("Worker started")
	defer workerLogger.Warnf("Worker ded")

	for {
		job := <-JobQueue
		if err := job.init(workerLogger); err != nil {
			workerLogger.WithError(err).Errorf("Error initializing job")
			job.clean()
			continue
		}

		if err := worker.RunJob(job); err != nil {
			workerLogger.WithError(err).Errorf("Error running job")
			job.clean()
			continue
		}
		job.timings.doneReadingTime = time.Now()
		job.wg.Wait() // Wait for container to exit.
		job.timings.completedTime = time.Now()

		job.SendOutput()
		job.clean()
	}
}

func (worker *Worker) RunJob(job *Job) error {
	logs, err := job.StartJob()
	if err != nil {
		return errors.Wrap(err, "Error starting job")
	}
	go job.waitAndKill(logs)

	// Read the output from the container.
	job.stdOut, job.stdErr, err = docker.ReadLogOutputs(logs)
	if err != nil {
		job.chans.errCh <- errors.Wrap(err, "Error reading container output")
	} else {
		job.chans.doneCh <- true
	}
	job.timings.doneReadingTime = time.Now()

	return nil
}
