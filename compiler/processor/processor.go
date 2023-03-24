package processor

import (
	"context"
	"sync"
	"time"

	"github.com/Squwid/bytegolf/lib/api"
	"github.com/Squwid/bytegolf/lib/sqldb"
	"github.com/sirupsen/logrus"
)

// benchmarkTestMultiplier is the number of times the benchmark test case should run.
const benchmarkTestMultiplier = 30

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
}

func NewSubmission(id string) *Submission {
	return &Submission{
		ID:         id,
		ReceivedAt: time.Now(),
		logger:     logrus.WithField("SubID", id),
	}
}

func ProcessMessage(ctx context.Context, id string) {
	sub := NewSubmission(id)

	// TODO: Handle errors here where they can report upstream.
	s, err := api.GetSubmission(ctx, sub.ID)
	if err != nil {
		sub.logger.WithError(err).Errorf("Error getting submission")
		return
	}
	if s == nil {
		sub.logger.Warnf("Submission not found")
		return
	}
	sub.Submission = s

	hole, err := api.GetHole(ctx, sub.Submission.Hole)
	if err != nil {
		sub.logger.WithError(err).Errorf("Error getting hole")
		return
	}
	sub.Hole = hole

	if err := api.UpdateSubmissionStatus(ctx, sub.ID, 1); err != nil {
		sub.logger.WithError(err).Errorf("Error updating submission status")
	} else {
		sub.logger.Debugf("Updated submission status to 1 (RUNNING)")
	}

	// Generate and send jobs to workers.
	sub.GenerateJobs()
	for _, job := range sub.jobs {
		JobQueue <- job
	}
	go sub.waitAndWriteToDB()
	sub.jobWG.Wait()
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
		job.logger.Debugf("Writing job output (%v %v) to DB",
			job.Output.SubmissionID,
			job.Output.TestID)

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

// GenerateJobs populates the counts of jobs and the waitgroup.
func (sub *Submission) GenerateJobs() {
	for _, test := range sub.Hole.TestsDB {
		var testCount = 1
		if test.Benchmark {
			testCount = benchmarkTestMultiplier
		}

		for i := 0; i < testCount; i++ {
			var job = NewJob(sub, test)
			sub.jobs = append(sub.jobs, job)
		}
	}

	sub.jobOutputs = make(chan *Job, len(sub.jobs))
	sub.jobWG = &sync.WaitGroup{}
	sub.jobWG.Add(len(sub.jobs))
}
