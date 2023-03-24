package main

import (
	"os"

	"github.com/Squwid/bytegolf/compiler/processor"
	"github.com/Squwid/bytegolf/lib/comms"
	"github.com/Squwid/bytegolf/lib/sqldb"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := sqldb.Open(); err != nil {
		logrus.WithError(err).Fatalf("Error connecting to db")
	}
	defer func() {
		if err := sqldb.Close(); err != nil {
			logrus.WithError(err).Errorf("")
		}
	}()

	comms.InitReceiver(os.Getenv("BG_USE_PUBSUB") == "true")
	comms.ReceiverImpl.Listen(processor.ProcessMessage)
}
