package github

import (
	"github.com/Squwid/bytegolf/secrets"
)

var githubClient *secrets.Client

func init() {
	githubClient = secrets.Must(secrets.GetClient("BGGH")).(*secrets.Client)
}

type client struct {
	ID     string `json:"client_id"`
	Secret string `json:"client_secret"`
}

func getClient() client {
	return client{
		ID:     githubClient.Client,
		Secret: githubClient.Secret,
	}
}
