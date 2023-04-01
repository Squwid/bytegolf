package main

import (
	"os"

	"github.com/Squwid/bytegolf/compiler/processor"
	"github.com/Squwid/bytegolf/lib/comms"
	"github.com/Squwid/bytegolf/lib/docker"
	"github.com/Squwid/bytegolf/lib/log"
	"github.com/Squwid/bytegolf/lib/sqldb"
)

func main() {
	if err := sqldb.Open(true); err != nil {
		log.GetLogger().WithError(err).Fatalf("Error connecting to db")
	}
	defer func() {
		if err := sqldb.Close(); err != nil {
			log.GetLogger().WithError(err).Errorf("")
		}
	}()
	docker.Init()

	if err := comms.InitReceiver(
		os.Getenv("BG_USE_PUBSUB") == "true"); err != nil {
		log.GetLogger().WithError(err).Fatalf("Error initializing receiver")
	}

	comms.ReceiverImpl.Listen(processor.ProcessMessage)
}
