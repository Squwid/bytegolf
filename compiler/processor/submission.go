package processor

import (
	"context"
	"fmt"
	"os"
	"sync"
	"text/template"
	"time"

	"github.com/Squwid/bytegolf/lib/api"
	"github.com/Squwid/bytegolf/lib/log"
	"github.com/Squwid/bytegolf/lib/sqldb"
	"github.com/sirupsen/logrus"
)

type Submission struct {
	ID         string
	ReceivedAt time.Time
	Submission *api.SubmissionDB
	Hole       *api.HoleDB

	logger *logrus.Entry

	// Job details.
	jobs       []*Job
	jobWG      *sync.WaitGroup
	jobOutputs chan *Job
	tempDir    string // Directory where the temp code file is stored on the host.
}

func NewSubmission(id string) *Submission {
	return &Submission{
		ID:         id,
		ReceivedAt: time.Now(),
		logger:     log.GetLogger().WithField("SubID", id),
	}
}

// Initialize the submission by creating the temp directory and writing
// all input files to it.
func (sub *Submission) Init() error {
	dir, err := os.MkdirTemp("", fmt.Sprintf("bg-%s_", sub.ID))
	if err != nil {
		return err
	}
	sub.tempDir = dir

	for _, test := range sub.Hole.TestsDB {
		tmpl, err := template.New(test.Name).Parse(test.Boilerplate)
		if err != nil {
			return err
		}

		f, err := os.Create(fmt.Sprintf("%s/main-%v.%s", sub.tempDir,
			test.ID, sub.Hole.LanguageDB.Extension))
		if err != nil {
			return err
		}

		if err := tmpl.Execute(f, struct{ UserSolution string }{
			UserSolution: sub.Submission.Script,
		}); err != nil {
			return err
		}
		f.Close()
	}

	return nil
}

func (sub *Submission) Clean() error {
	return os.RemoveAll(sub.tempDir)
}

// GenerateJobs populates the counts of jobs and the waitgroup.
func (sub *Submission) GenerateJobs() {

	// TODO: Setup a way to only setup one job
	for _, test := range sub.Hole.TestsDB {
		if !test.Benchmark {
			sub.jobs = append(sub.jobs, NewJob(sub, test))
			break
		}
	}

	// for _, test := range sub.Hole.TestsDB {
	// 	var testCount = 1
	// 	// if true {
	// 	if test.Benchmark {
	// 		testCount = benchmarkTestMultiplier
	// 	}

	// 	for i := 0; i < testCount; i++ {
	// 		var job = NewJob(sub, test)
	// 		sub.jobs = append(sub.jobs, job)
	// 	}
	// }

	sub.jobOutputs = make(chan *Job, len(sub.jobs))
	sub.jobWG = &sync.WaitGroup{}
	sub.jobWG.Add(len(sub.jobs))
}

// waitAndWriteToDB waits for jobs to finish and writes them as they arrive
// to the database.
func (sub *Submission) waitAndWriteToDB() {
	for job := range sub.jobOutputs {
		// TODO: Now this makes the job "correct" but there should be another enum
		// option for failure during check.
		if err := job.checkCorrectness(); err != nil {
			sub.logger.WithError(err).Errorf("Error checking correctness")
		}

		job.Output.Correct = job.correct
		job.logger.Debugf("Writing job output to DB")

		if _, err := sqldb.DB.
			NewInsert().
			Model(job.Output).
			Exec(context.Background()); err != nil {
			job.logger.WithError(err).Errorf("Error inserting job output")
		}
		sub.jobWG.Done()
	}
	close(sub.jobOutputs)
}
