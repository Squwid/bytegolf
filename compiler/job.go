package main

import (
	"fmt"
	"os"

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

	// Internal job details
	dir         string
	file        string
	containerID string

	outputCh  chan *Job
	errCh     chan error
	timeoutCh chan bool

	logger *logrus.Entry
	output *api.JobOutputDB
}

func NewJob(sub *api.SubmissionDB, hole *api.HoleDB,
	test *api.TestDB, jobOutputs chan *Job) *Job {
	return &Job{
		ID:         api.RandomString(10),
		Language:   hole.LanguageDB,
		Test:       test,
		Submission: sub,
		outputCh:   jobOutputs,
		errCh:      make(chan error, 1),
		timeoutCh:  make(chan bool, 1),
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
