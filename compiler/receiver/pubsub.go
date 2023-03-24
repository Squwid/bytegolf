package receiver

import (
	"context"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/Squwid/bytegolf/compiler/processor"
	"github.com/sirupsen/logrus"
)

const subName = "bgcompiler-rpi-local-sub"

// PubSub is the type that implements the Receiver interface for the pub/sub.
// This is used in production and is skippable if the http receiver is used instead.
type PubSub struct {
	client *pubsub.Client
	sub    *pubsub.Subscription
}

func (p *PubSub) Init() error {
	logrus.Infof("Initializing PubSub receiver (sub: %s)", subName)

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, os.Getenv("GCP_PROJECT_ID"))
	if err != nil {
		return err
	}
	p.client = client
	p.sub = p.client.Subscription(subName)
	return nil
}

func (p *PubSub) Start() {
	go func() {
		for {
			if err := p.sub.Receive(context.Background(), pubsubHandler); err != nil {
				panic(err)
			}
		}
	}()
}

func pubsubHandler(ctx context.Context, msg *pubsub.Message) {
	// TODO: Update the submission database with a note that the compiler is working on it.
	// TODO: Only ack if there is room for the submission in the queue.
	logger := logrus.WithFields(logrus.Fields{
		"Action":       "Handler",
		"MessageID":    msg.ID,
		"SubmissionID": string(msg.Data),
	})
	logger.Infof("Recieved message")

	msg.Ack()
	processor.ProcessMessage(ctx, string(msg.Data))
}
