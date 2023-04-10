package processor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Squwid/bytegolf/lib/api"
	"github.com/Squwid/bytegolf/lib/log"
	"github.com/Squwid/bytegolf/lib/sqldb"
	"github.com/sirupsen/logrus"
)

// benchmarkTestMultiplier is the number of times the benchmark test case should run.
const benchmarkTestMultiplier = 10

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
		logger:     log.GetLogger().WithField("SubID", id),
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

	if err := api.UpdateSubmissionStatus(ctx, sub.ID, api.StatusRunning); err != nil {
		sub.logger.WithError(err).Errorf("Error updating submission status")
	} else {
		sub.logger.Debugf("Updated submission status to %v %s",
			api.StatusRunning, api.StatusRunningStr)
	}

	sub.GenerateJobs()
	sub.logger = sub.logger.WithField("Jobs", len(sub.jobs))
	for _, job := range sub.jobs {
		JobQueue <- job
	}

	go sub.waitAndWriteToDB()

	start := time.Now()
	sub.logger.Infof("Jobs sent to queue. Waiting for jobs to finish...")
	sub.jobWG.Wait()

	passed, countPassed := checkIfPassed(sub.jobs)
	avgs := calculateAverages(sub.jobs)
	if err := api.CompleteSubmission(ctx, sub.ID, passed, avgs); err != nil {
		sub.logger.WithError(err).Errorf("Error updating submission")
	} else {
		sub.logger.Debugf("Updated submission status to %v %s", api.StatusSuccess, api.StatusSuccessStr)
	}

	sub.logger.WithFields(logrus.Fields{
		"TotalMS":       time.Since(start).Milliseconds(),
		"AvgCPU":        avgs.AvgCPU,
		"AvgMem":        avgs.AvgMem,
		"Passed":        passed,
		"PercentPassed": fmt.Sprintf("%.2f%%", percent(countPassed, len(sub.jobs))),
	}).Infof("All jobs finished.")
}

func percent(countPassed, jobs int) float64 {
	if jobs == 0 {
		return 0
	}
	return (float64(countPassed) / float64(jobs)) * 100
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

// GenerateJobs populates the counts of jobs and the waitgroup.
func (sub *Submission) GenerateJobs() {

	// TODO: Setup a way to only setup one job
	// for _, test := range sub.Hole.TestsDB {
	// 	if !test.Benchmark {
	// 		sub.jobs = append(sub.jobs, NewJob(sub, test))
	// 		break
	// 	}
	// }

	for _, test := range sub.Hole.TestsDB {
		var testCount = 1
		// if true {
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
