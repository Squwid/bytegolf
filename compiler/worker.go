package main

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/Squwid/bytegolf/lib/api"
	"github.com/Squwid/bytegolf/lib/docker"
	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const timeout = 5 * time.Second
const workerCount = 2
const jobBacklog = 5000
const bytesToRead = 4096

var jobQueue = make(chan *Job, jobBacklog)

type Worker struct {
	ID string

	lock *sync.Mutex
}

type Stats struct {
	CPU int64
	Mem int64
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
		ctx := context.Background()
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

		// Read stats as they are available.
		stats, err := docker.Client.Stats(ctx, containerID)
		if err != nil {
			job.logger.WithError(err).Error("Failed to get container stats")
			job.errCh <- errors.Wrap(err, "Failed to get container stats")
			job.clean()
			continue
		}

		go waitAndGetContainerStats(stats.Body, job)
		go waitAndKillContainer(ctx, reader, job)

		// Read the output from the container.
		out, err := readAmount(reader, bytesToRead, job.logger)
		if err != nil {
			job.logger.WithError(err).Errorf("Error reading bytes")
			job.errCh <- errors.Wrap(err, "Error reading bytes")
			job.clean()
			continue
		}
		job.doneCh <- true // Done reading output.
		dur := time.Since(start)
		job.wg.Wait() // Wait for job to be "completed"

		var jobOut = &api.JobOutputDB{
			StdOut:   string(out),
			Duration: dur.Milliseconds(),
			Memory:   int64(job.stats.Mem),
			CPU:      int64(job.stats.CPU),
			ExitCode: 0, // TODO: Populate exit code.
		}
		job.logger.Debugf("Finished job in %vms", dur.Milliseconds())
		job.output = jobOut
		job.outputCh <- job
		job.clean()
	}
}

func waitAndGetContainerStats(reader io.ReadCloser, job *Job) {
	defer job.wg.Done()
	defer reader.Close()

	dec := json.NewDecoder(reader)
	var stats = &Stats{}
	for {
		var v *types.StatsJSON
		if err := dec.Decode(&v); err != nil {
			job.stats = stats
			if err == io.EOF {
				return
			}
			job.logger.WithError(err).Errorf("Error decoding stats")
			return
		}

		if stats.CPU < int64(v.CPUStats.CPUUsage.TotalUsage) {
			stats.CPU = int64(v.CPUStats.CPUUsage.TotalUsage)
		}
		if stats.Mem < int64(v.MemoryStats.Usage) {
			stats.Mem = int64(v.MemoryStats.Usage)
		}
	}
}

func waitAndKillContainer(ctx context.Context, reader io.ReadCloser, job *Job) {
	defer job.wg.Done()

	// Wait for the job to finish or timeout.
	select {
	case <-job.doneCh:
		job.logger.Debugf("Got done reading signal to close reader")
	case <-time.After(timeout):
		job.timedOut = true
		job.logger.Debugf("Job timed out")
	case <-job.errCh:
		job.logger.Debugf("Got error signal to close reader")
	}

	if err := reader.Close(); err != nil {
		job.logger.WithError(err).Debugf("Error closing reader")
	}
	if err := docker.Client.Kill(ctx, job.containerID); err != nil {
		job.logger.WithError(err).Debugf("Error killing container")
	}

	if err := docker.Client.Delete(ctx, job.containerID); err != nil {
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
