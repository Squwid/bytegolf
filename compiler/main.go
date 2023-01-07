package main

import (
	"context"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/Squwid/bytegolf/lib/api"
	"github.com/Squwid/bytegolf/lib/sqldb"
	"github.com/sirupsen/logrus"
)

const timeout = 5 * time.Second

// const subName = "bgcompiler-sub"
const subName = "bgcompiler-rpi-local-sub"

// testRepeatMultiplier is the number of times to repeat the test case.
const testRepeatMultiplier = 20

func main() {
	if err := sqldb.Open(); err != nil {
		logrus.WithError(err).Fatalf("Error connecting to db")
	}
	defer func() {
		if err := sqldb.Close(); err != nil {
			logrus.WithError(err).Errorf("")
		}
	}()

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, os.Getenv("GCP_PROJECT_ID"))
	if err != nil {
		logrus.WithError(err).Fatalf("Error creating pub/sub subscriber")
	}
	sub := client.Subscription(subName)

	for {
		if err = sub.Receive(context.Background(), handler); err != nil {
			logrus.WithError(err).Fatalf("Error recieving message\n")
		}
	}
}

// CompiledSubmission is the type that holds all of the individual test results
// for a submission.
type CompiledSubmission struct {
	jobOutputs chan *Job
	wg         *sync.WaitGroup
	jobs       []*Job
}

func newCompiledSubmission(jobs int) *CompiledSubmission {
	wg := &sync.WaitGroup{}
	wg.Add(jobs)
	return &CompiledSubmission{
		jobOutputs: make(chan *Job, jobs),
		wg:         wg,
	}
}

func handler(ctx context.Context, m *pubsub.Message) {
	start := time.Now()
	// TODO: Update the submission database with a note that the compiler is working on it.
	// TODO: Only ack if there is room for the submission in the queue.
	logger := logrus.WithFields(logrus.Fields{
		"Action":       "Handler",
		"MessageID":    m.ID,
		"SubmissionID": string(m.Data),
	})
	logger.Infof("Recieved message")

	// Parse the object coming from the message queue. Fetch test cases from
	// the database and send each test as a seperate entity to the workers.
	sub, err := api.GetSubmission(ctx, string(m.Data))
	if err != nil {
		logger.WithError(err).Errorf("Error getting submission")
		return
	}

	hole, err := api.GetHole(ctx, sub.Hole)
	if err != nil {
		logger.WithError(err).Errorf("Error getting hole")
		return
	}
	m.Ack()
	logger.Infof("Acked message")

	cs := newCompiledSubmission(len(hole.TestsDB) * testRepeatMultiplier)

	for _, test := range hole.TestsDB {
		for i := 0; i < testRepeatMultiplier; i++ {
			var job = NewJob(sub, hole, test, cs.jobOutputs)
			jobQueue <- job
		}
	}

	go func(cs *CompiledSubmission, logger *logrus.Entry) {
		for job := range cs.jobOutputs {
			job.output.SubmissionID = job.Submission.ID
			job.output.TestID = job.Test.ID
			logger.Infof("Writing job output (%v %v) to DB", job.output.SubmissionID,
				job.output.TestID)
			cs.jobs = append(cs.jobs, job)

			if _, err := sqldb.DB.NewInsert().Model(job.output).Exec(ctx); err != nil {
				logger.WithError(err).Errorf("Error inserting job output")
			}
			cs.wg.Done()
		}
		close(cs.jobOutputs)
	}(cs, logger)

	cs.wg.Wait()

	// TODO: Update the submission database with the results of the tests.
	logger.Infof("DONE AFTER %vms", time.Since(start).Milliseconds())
}
