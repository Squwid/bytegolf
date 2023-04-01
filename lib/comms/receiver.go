package comms

import (
	"context"

	"github.com/Squwid/bytegolf/lib/log"
)

const (
	CommsTypeHttp   = "Http"
	CommsTypePubSub = "PubSub"
)

type Receiver interface {
	// Init should initialize the receiver and open any necessary connections.
	Init() error

	// Start should start the receiver in the main thread and stay alive.
	Listen(func(context.Context, string))
}

var ReceiverImpl Receiver = nil

// InitReceiver initializes the receiver implementation. If usePubSub is true, it will
// use Google Cloud Pub/Sub. Otherwise, it will use HTTP.
func InitReceiver(usePubSub bool) error {
	var rt = CommsTypePubSub
	if !usePubSub {
		rt = CommsTypeHttp
	}

	log.GetLogger().WithField("ReceiverType", rt).
		Infof("Initializing receiver...")

	if usePubSub {
		ReceiverImpl = &PubSub{}
	} else {
		ReceiverImpl = &Http{}
	}

	return ReceiverImpl.Init()
}
