package processor

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/Squwid/bytegolf/lib/docker"
	"github.com/Squwid/bytegolf/lib/log"
)

const (
	timeout = 10 * time.Second

	workerCount = 4
	jobBacklog  = 5000
	bytesToRead = 4096
)

var JobQueue = make(chan *Job, jobBacklog)

type Worker interface {
	Start()
}

type WorkerData struct {
	ID string

	lock *sync.Mutex
}

var WorkerPool [workerCount]Worker

func init() {
	for i := 0; i < workerCount; i++ {
		worker := NewWorker(i)
		WorkerPool[i] = worker
		go worker.Start()
	}
}

func NewWorker(id int) *WorkerData {
	return &WorkerData{
		ID:   fmt.Sprintf("w%v", id),
		lock: &sync.Mutex{},
	}
}

func (worker *WorkerData) Start() {
	wl := log.GetLogger().WithField("Worker", worker.ID)
	wl.Info("Worker started")
	defer wl.Warnf("Worker ded")

	for {
		job := <-JobQueue
		if err := job.init(wl); err != nil {
			wl.WithError(err).Errorf("Error initializing job")
			continue
		}

		containerID, reader, err := job.createContainer()
		if err != nil {
			job.logAndReportError(err, "Error creating container")
			continue
		}
		job.containerID = containerID

		// Continuously poll for stats.
		stats, err := docker.Client.Stats(job.ctx, job.containerID)
		if err != nil {
			job.logAndReportError(err, "Error getting container stats")
			continue
		}

		go waitAndGetContainerStats(stats.Body, job)
		go waitAndKillContainer(job.ctx, reader, job)

		// Read the output from the container.
		stdOut, stdErr, err := docker.ReadLogOutputs(reader)
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

func waitAndKillContainer(ctx context.Context, reader io.ReadCloser, job *Job) {
	defer job.wg.Done()

	// Wait for the job to finish or timeout.
	select {
	case <-job.chans.doneCh:
		job.logger.Debugf("Got done reading signal to close reader")
	case <-time.After(timeout):
		job.timedOut = true
		job.logger.Debugf("Job timed out")
	case <-job.chans.errCh:
		job.logger.Debugf("Got error signal to close reader")
	}

	_ = reader.Close()
	_ = docker.Client.Kill(ctx, job.containerID)
	_ = docker.Client.Delete(ctx, job.containerID)
}
