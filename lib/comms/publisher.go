package comms

import (
	"context"

	"github.com/Squwid/bytegolf/lib/log"
)

type Publisher interface {
	// Init should initialize the publisher and open any necessary connections.
	Init() error

	// Publish should send a message to the compiler.
	Publish(ctx context.Context, message []byte) error
}

var PublisherImpl Publisher = nil

// InitPublisher initializes the publisher implementation. If usePubSub is true, it will
// use Google Cloud Pub/Sub. Otherwise, it will use HTTP.
func InitPublisher(usePubSub bool) error {
	var pt = CommsTypePubSub
	if !usePubSub {
		pt = CommsTypeHttp
	}

	log.GetLogger().WithField("PublisherType", pt).
		Infof("Initializing compiler publisher...")

	if usePubSub {
		PublisherImpl = &PubSub{}
	} else {
		PublisherImpl = &Http{}
	}

	return PublisherImpl.Init()
}
