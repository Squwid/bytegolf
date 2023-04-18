package processor

import (
	"context"
	"fmt"
	"time"

	"github.com/Squwid/bytegolf/lib/api"
	"github.com/sirupsen/logrus"
)

// benchmarkTestMultiplier is the number of times the benchmark test case should run.
const benchmarkTestMultiplier = 15

func ProcessMessage(ctx context.Context, id string) {
	// TODO: If any of these initialization processes fail, fail the request
	// and dont accept the message from pubsub.
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

	if sub.Init() != nil {
		sub.logger.Errorf("Error initializing submission")
		return
	}

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
