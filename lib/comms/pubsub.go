package comms

import (
	"context"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/sirupsen/logrus"
)

const subName = "bgcompiler-rpi-local-sub"
const topic = "testing"

// PubSub is the type that implements the Receiver interface for the pub/sub.
// This is used in production and is skippable if the http receiver is used instead.
type PubSub struct {
	client *pubsub.Client
	sub    *pubsub.Subscription
	topic  *pubsub.Topic
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
	p.topic = p.client.Topic(topic)
	return nil
}

func (p *PubSub) Publish(ctx context.Context, message []byte) error {
	res := p.topic.Publish(ctx, &pubsub.Message{
		Data: message,
	})
	_, err := res.Get(ctx)
	return err

}

func (p *PubSub) Listen(processor func(context.Context, string)) {
	for {
		if err := p.sub.Receive(context.Background(),
			pubsubHandler(processor)); err != nil {
			panic(err)
		}
	}
}

func pubsubHandler(processor func(context.Context, string)) func(context.Context, *pubsub.Message) {
	return func(ctx context.Context, msg *pubsub.Message) {
		// TODO: Only ack if there is room for the submission in the queue.
		logger := logrus.WithFields(logrus.Fields{
			"Action":       "Handler",
			"MessageID":    msg.ID,
			"SubmissionID": string(msg.Data),
		})
		logger.Infof("Recieved message")

		msg.Ack()
		processor(ctx, string(msg.Data))
	}
}
