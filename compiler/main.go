package main

import (
	"context"
	"fmt"
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

// benchmarkTestMultiplier is the number of times the benchmark test case should run.
const benchmarkTestMultiplier = 30

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
		jobs:       []*Job{},
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
	logger.Debugf("Acked message")

	cs := newCompiledSubmission(len(hole.TestsDB) + (benchmarkTestMultiplier - 1))

	// Create a job for each test. Create 20x jobs for the benchmark test case.
	for _, test := range hole.TestsDB {
		var testCount = 1
		if test.Benchmark {
			testCount = benchmarkTestMultiplier
		}

		for i := 0; i < testCount; i++ {
			var job = NewJob(sub, hole, test, cs.jobOutputs)
			cs.jobs = append(cs.jobs, job)
			jobQueue <- job
		}
	}

	// Wait for all jobs to finish running and writing to the DB.
	go waitAndWriteToDB(ctx, cs, logger)
	cs.wg.Wait()

	// Check if all test cases passed. Error should never occur
	// but better to have handler than a panic.
	// TODO: Rather than a single pass fail do one for each test case.
	// TODO: Compile 1 regex per test case rather than 1 regex per submission.
	var totalPassed = 0
	var passed = true
	for _, job := range cs.jobs {
		if !job.correct {
			passed = false
		} else {
			totalPassed++
		}
	}

	// Average CPU times.
	var cpuTimes = []int64{}
	for i := range cs.jobs {
		if cs.jobs[i].Test.Benchmark {
			cpuTimes = append(cpuTimes, cs.jobs[i].output.Duration)
		}
	}
	var cpuAverage = average(cpuTimes)

	// TODO: Update the submission database with the results of the tests.
	logger.WithFields(logrus.Fields{
		"Jobs":          len(cs.jobs),
		"TotalMS":       time.Since(start).Milliseconds(),
		"BenchmarkCPU":  cpuAverage,
		"Passed":        passed,
		"PercentPassed": fmt.Sprintf("%.2f%%", (float64(totalPassed)/float64(len(cs.jobs)))*100),
	}).Infof("Finished submission")
}

func waitAndWriteToDB(ctx context.Context, cs *CompiledSubmission, logger *logrus.Entry) {
	for job := range cs.jobOutputs {
		if err := job.checkCorrectness(); err != nil {
			logger.WithError(err).Errorf("Error checking correctness")
		}

		job.output.SubmissionID = job.Submission.ID
		job.output.TestID = job.Test.ID
		job.output.Correct = job.correct
		logger.Debugf("Writing job output (%v %v) to DB", job.output.SubmissionID,
			job.output.TestID)

		if _, err := sqldb.DB.NewInsert().Model(job.output).Exec(ctx); err != nil {
			logger.WithError(err).Errorf("Error inserting job output")
		}
		cs.wg.Done()
	}
	close(cs.jobOutputs)
}

func average(values []int64) int64 {
	var sum int64
	for _, v := range values {
		sum += v
	}
	return sum / int64(len(values))
}
