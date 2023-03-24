package receiver

import "os"

type Receiver interface {
	// Init should initialize the receiver and open any necessary connections.
	Init() error

	// Start should start the receiver in a goroutine, and keep alive.
	Start()
}

var ReceiverImpl Receiver

func init() {
	if os.Getenv("USE_PUBSUB") == "true" {
		ReceiverImpl = &PubSub{}
	} else {
		ReceiverImpl = &Http{}
	}
}
