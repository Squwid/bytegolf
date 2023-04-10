package processor

import (
	"fmt"
	"sync"
	"time"

	"github.com/Squwid/bytegolf/lib/docker"
	"github.com/Squwid/bytegolf/lib/log"
)

const (
	timeout = 10 * time.Second

	jobBacklog  = 5000
	bytesToRead = 4096
	workerCount = 8
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
	wl := log.GetLogger().WithField("Worker", worker.ID)
	wl.Info("Worker started")
	defer wl.Warnf("Worker ded")

	for {
		job := <-JobQueue
		if err := job.init(wl); err != nil {
			job.logAndReportError(err, "Error initializing job")
			continue
		}

		logs, err := job.StartJob()
		if err != nil {
			job.logAndReportError(err, "Error starting job")
			continue
		}
		go job.waitAndKill(logs)

		// Read the output from the container.
		stdOut, stdErr, err := docker.ReadLogOutputs(logs)
		if err != nil {
			job.logAndReportError(err, "Error reading container output")
			continue
		}

		job.chans.doneCh <- true
		job.timings.doneReadingTime = time.Now()
		job.wg.Wait()
		job.timings.completedTime = time.Now()
		job.SetOutput(
			job.timings.doneReadingTime.Sub(job.timings.initTime),
			string(stdOut.Output()),
			string(stdErr.Output()),
		)
		job.Sub.jobOutputs <- job // Send job output to submission.

		job.clean()
	}
}
