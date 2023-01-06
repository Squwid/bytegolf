package sqldb

import (
	"context"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/sirupsen/logrus"
)

var PubSubClient *PubSub

type PubSub struct {
	client *pubsub.Client
}

var topic = os.Getenv("PUBSUB_TOPIC")

func init() {
	client, err := NewPubSub(context.Background())
	if err != nil {
		logrus.WithError(err).Fatal("Error creating pubsub client")
	}
	PubSubClient = client
}

func NewPubSub(ctx context.Context) (*PubSub, error) {
	client, err := pubsub.NewClient(ctx, os.Getenv("PUBSUB_PROJECT"))
	if err != nil {
		return nil, err
	}
	logrus.Infof("PubSub client created")
	return &PubSub{
		client: client,
	}, nil
}

func (ps *PubSub) Publish(ctx context.Context, message []byte) error {
	t := ps.client.Topic(topic)
	result := t.Publish(ctx, &pubsub.Message{
		Data: message,
	})
	_, err := result.Get(ctx)
	return err
}
