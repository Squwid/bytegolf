package processor

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/Squwid/bytegolf/lib/api"
	"github.com/Squwid/bytegolf/lib/docker"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Job struct {
	// TODO: This doesnt match the autoincrementing ID in the DB.
	ID     string
	Sub    *Submission
	Test   *api.TestDB
	Stats  *Stats
	Output *api.JobOutputDB

	correct  bool // Correct based on the test regex.
	timedOut bool // True if the job timed out during execution.

	timings JobTimings

	dir              string // Directory where the temp code file is stored on the host.
	absoluteFilePath string // Absolute path to the code file on the host.
	fileName         string // Name of the file, with the extension.
	containerID      string

	chans JobChannels

	// wg waits for stats and container removal to
	// finish before the job is considered done.
	wg     *sync.WaitGroup
	logger *logrus.Entry
	ctx    context.Context
}

// JobChannels holds the channels used to communicate with the job.
type JobChannels struct {
	errCh  chan error
	doneCh chan bool // Signal that reading input is done.
}

// JobTimings is used to track the time it takes for a job to be
// processed from different stages.
type JobTimings struct {
	// createdTime is when the job is created by the processor.
	createdTime time.Time

	// initTime is when the job is initialized by the worker.
	initTime time.Time

	// doneReadingTime is when the job is done reading output
	// from the container.
	doneReadingTime time.Time

	// completedTime is when the job is done waiting for container
	// to be cleaned up.
	completedTime time.Time
}

func NewJob(sub *Submission, test *api.TestDB) *Job {
	wg := &sync.WaitGroup{}
	wg.Add(2) // 1 for containerStats, 1 for waitAndKillContainer
	return &Job{
		ID:      api.RandomString(10),
		Sub:     sub,
		Test:    test,
		correct: true,
		wg:      wg,
		chans: JobChannels{
			errCh:  make(chan error, 1),
			doneCh: make(chan bool, 1),
		},
		timings: JobTimings{
			createdTime: time.Now(),
		},
	}
}

func (job *Job) SetOutput(dur time.Duration, stdOut, stdErr string) {
	job.Output = &api.JobOutputDB{
		SubmissionID: job.Sub.ID,
		TestID:       job.Test.ID,

		StdOut:   stdOut,
		StdErr:   stdErr,
		Duration: dur.Milliseconds(),
		Memory:   int64(job.Stats.Mem),
		CPU:      int64(job.Stats.CPU),
		ExitCode: 0, // TODO: Populate exit code.
	}
}

func (job *Job) init(workerLogger *logrus.Entry) error {
	job.logger = workerLogger.WithFields(logrus.Fields{
		"JobID":  job.ID,
		"SubID":  job.Sub.ID,
		"TestID": job.Test.ID,
	})
	job.ctx = context.Background()
	job.timings.initTime = time.Now()
	return job.writeFiles()
}

func (job *Job) writeFiles() error {
	dir, err := os.MkdirTemp("", "bg")
	if err != nil {
		return err
	}
	job.dir = dir
	job.fileName = fmt.Sprintf("main.%s",
		job.Sub.Hole.LanguageDB.Extension)
	job.absoluteFilePath = fmt.Sprintf("%s/%s", job.dir,
		job.fileName)

	return os.WriteFile(job.absoluteFilePath,
		[]byte(job.Sub.Submission.Script), 0755)
}

// createContainer is a wrapper around docker.Client.Create to put logic for
// job creation in a single place.
func (job *Job) createContainer() (string, io.ReadCloser, error) {
	return docker.Client.Create(
		job.Sub.Hole.LanguageDB.Image,
		job.absoluteFilePath,
		job.fileName,
		job.Sub.Hole.LanguageDB.Cmd,
		job.ID,
		fmt.Sprintf("%s/%s", job.Test.Hole, job.Test.Input),
		job.logger,
	)
}

func (job *Job) logAndReportError(err error, msg string) {
	job.logger.WithError(err).Error(msg)
	job.chans.errCh <- errors.Wrap(err, msg)
	job.clean()
}

func (job *Job) clean() error {
	return os.RemoveAll(job.dir)
}

func (job *Job) checkCorrectness() error {
	regex, err := regexp.Compile(job.Test.OutputRegex)
	if err != nil {
		return err
	}
	if !regex.MatchString(job.Output.StdOut) {
		job.correct = false
	}
	return nil
}
