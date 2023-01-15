package main

import (
	"context"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/Squwid/bytegolf/lib/api"
	"github.com/Squwid/bytegolf/lib/docker"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const timeout = 5 * time.Second
const workerCount = 4
const jobBacklog = 5000
const bytesToRead = 4096

var jobQueue = make(chan *Job, jobBacklog)

type Worker struct {
	ID string

	lock *sync.Mutex
}

func init() {
	for i := 0; i < workerCount; i++ {
		worker := newWorker()
		go worker.Start()
	}
}

func newWorker() *Worker {
	return &Worker{
		ID:   api.RandomString(4),
		lock: &sync.Mutex{},
	}
}

func (worker *Worker) Start() {
	workerLogger := logrus.WithField("Worker", worker.ID)
	workerLogger.Info("Worker started")
	defer workerLogger.Warnf("Worker ded")

	for {
		job := <-jobQueue
		job.logger = workerLogger.WithField("JobID", job.ID)
		job.logger.Debugf("Job started (%s, %v)", job.Submission.ID, job.Test.ID)

		if err := job.init(); err != nil {
			job.logger.WithError(err).Errorf("Error initializing job")
			continue
		}

		// Run the job.
		start := time.Now()
		containerID, reader, err := docker.Client.Create(
			job.Language.Image,
			job.dir,
			job.Language.Cmd,
			job.file,
			job.ID,
			job.logger)
		if err != nil {
			job.logger.WithError(err).Error("Failed to create container")
			job.errCh <- err
			job.clean()
			continue
		}
		job.containerID = containerID

		// Signal that we are done reading for the container.
		doneChan := make(chan bool, 1)
		go waitAndKillContainer(doneChan, reader, job)

		// Read the output from the container.
		out, err := readAmount(reader, bytesToRead, job.logger)
		if err != nil {
			job.logger.WithError(err).Errorf("Error reading bytes")
			job.errCh <- errors.Wrap(err, "Error reading bytes")
			job.clean()
			continue
		}
		doneChan <- true // Done reading output.

		dur := time.Since(start)
		ms := dur.Milliseconds()

		var jobOut = &api.JobOutputDB{
			StdOut:   string(out),
			Duration: ms,
			ExitCode: 0, // TODO: Populate exit code.
		}
		job.logger.Debugf("Finished job in %vms", ms)
		job.output = jobOut
		job.outputCh <- job
		job.clean()
	}
}

func waitAndKillContainer(doneChan chan bool, reader io.ReadCloser, job *Job) {
	select {
	case <-doneChan:
		job.logger.Debugf("Got done reading signal to close reader")
	case <-time.After(timeout):
		job.timedOut = true
		job.logger.Debugf("Job timed out")
	case <-job.errCh:
		job.logger.Debugf("Got error signal to close reader")
	}

	// TODO: Come up with a better way to kill containers.
	if err := reader.Close(); err != nil {
		job.logger.WithError(err).Debugf("Error closing reader")
	}
	if err := docker.Client.Kill(context.Background(), job.containerID); err != nil {
		job.logger.WithError(err).Debugf("Error killing container")
	}
	if err := docker.Client.Delete(context.Background(), job.containerID); err != nil {
		job.logger.WithError(err).Debugf("Error deleting container")
	}
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

		// Read everything we need
		if read >= amount {
			break
		}
	}

	return bytes[0:read], nil
}
