package auth

import (
	"github.com/Squwid/bytegolf/secrets"
	log "github.com/sirupsen/logrus"
)

const cookieName = "bg-token"

var id, secret, state string
var jwtKey []byte

func init() {
	// Get the github client info
	ghClient, err := secrets.GetClient("BGGH")
	if err != nil {
		log.Errorf("Error getting Github secret: %v")
		panic(err)
	}

	// Get the state
	stateClient, err := secrets.GetClient("BGSTATE")
	if err != nil {
		log.Errorf("Error getting state: %v", err)
		panic(err)
	}

	id = ghClient.Client
	secret = ghClient.Secret
	state = stateClient.Secret

	log.Infof("Successfully initialized github things")
}
