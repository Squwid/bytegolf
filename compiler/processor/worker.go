package processor

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/Squwid/bytegolf/lib/docker"
	"github.com/sirupsen/logrus"
)

const (
	timeout = 5 * time.Second

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
	wl := logrus.WithField("Worker", worker.ID)
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
		output, err := readAmount(reader, bytesToRead, job.logger)
		if err != nil {
			job.logAndReportError(err, "Error reading container output")
			continue
		}

		job.chans.doneCh <- true
		job.timings.doneReadingTime = time.Now()
		job.wg.Wait()
		job.timings.completedTime = time.Now()
		job.SetOutput(job.timings.doneReadingTime.Sub(job.timings.initTime), string(output))
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

	// TODO: Better error handle the closing of the docker container here.
	if err := reader.Close(); err != nil {
		job.logger.WithError(err).Debugf("Error closing reader")
	}

	docker.Client.Kill(ctx, job.containerID)
	docker.Client.Delete(ctx, job.containerID)
}

// Read only a certain amount of output without using a ton of memory.
func readAmount(r io.Reader, amount int, logger *logrus.Entry) ([]byte, error) {
	var bytes = make([]byte, amount) // buffer that gets returned.
	const readAmount = 1024          // Amount of bytes to read at a time.

	var i, read int
	for {
		bs := make([]byte, readAmount)
		n, err := r.Read(bs)

		if read < amount {
			read += copy(bytes[read:], bs[:n])
		}
		// This is what docker does in terms of checking if closed network connection.
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(),
				"use of closed network connection") {
				break
			}
			return nil, err
		}
		i++

		if read >= amount {
			break
		}
	}

	return bytes[0:read], nil
}
