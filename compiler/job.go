package main

import (
	"fmt"
	"os"
	"regexp"
	"sync"

	"github.com/Squwid/bytegolf/lib/api"
	"github.com/sirupsen/logrus"
)

type Job struct {
	// ID differs from SubmissionID since each test case run
	// needs a seperate ID.
	ID         string
	Language   *api.LanguageDB
	Test       *api.TestDB
	Submission *api.SubmissionDB

	correct  bool   // Correct based on the test regex.
	timedOut bool   // True if the job timed out during execution.
	stats    *Stats // CPU and memory usage

	// Internal job details
	dir         string
	file        string
	containerID string

	outputCh chan *Job
	errCh    chan error
	doneCh   chan bool // Signal that reading input is done.

	// wg waits for stats and container removal to
	// finish before the job is considered done.
	wg *sync.WaitGroup

	logger *logrus.Entry
	output *api.JobOutputDB
}

func NewJob(sub *api.SubmissionDB, hole *api.HoleDB,
	test *api.TestDB, jobOutputs chan *Job) *Job {
	wg := &sync.WaitGroup{}
	wg.Add(2) // 1 for containerStats, 1 for waitAndKillContainer
	return &Job{
		ID:         api.RandomString(10),
		Language:   hole.LanguageDB,
		Test:       test,
		Submission: sub,
		correct:    true,
		outputCh:   jobOutputs,
		errCh:      make(chan error, 1),
		doneCh:     make(chan bool, 1),
		wg:         wg,
	}
}

func (job *Job) init() error {
	return job.writeFiles()
}

func (job *Job) writeFiles() error {
	dir, err := os.MkdirTemp("", "bg")
	if err != nil {
		return err
	}
	job.dir = dir + "/"

	// Write code file.
	job.file = fmt.Sprintf("main.%s", job.Language.Extension)

	// Write input file.
	if job.Test.Input != "" {
		if err := os.WriteFile(job.dir+"input.txt",
			[]byte(job.Test.Input), 0644); err != nil {
			return err
		}
	}

	return os.WriteFile(job.dir+job.file,
		[]byte(job.Submission.Script), 0755)
}

func (job *Job) clean() error {
	return os.RemoveAll(job.dir)
}

func (job *Job) checkCorrectness() error {
	regex, err := regexp.Compile(job.Test.OutputRegex)
	if err != nil {
		return err
	}
	if !regex.MatchString(job.output.StdOut) {
		job.correct = false
	}
	return nil
}
