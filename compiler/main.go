package main

import (
	"context"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/Squwid/bytegolf/lib/api"
	"github.com/Squwid/bytegolf/lib/sqldb"
	"github.com/sirupsen/logrus"
)

const timeout = 5 * time.Second
const subName = "bgcompiler-sub"

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
			logrus.WithError(err).Errorf("Error recieving message\n")
		}
	}
}

func handler(ctx context.Context, m *pubsub.Message) {
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

	for _, test := range hole.TestsDB {
		var input = WorkerInput{
			Language:   hole.LanguageDB,
			Test:       test,
			Submission: sub,
		}
		workerQueue <- input
	}
}
