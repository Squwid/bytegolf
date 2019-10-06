package github

import "os"

type client struct {
	ID     string `json:"client_id"`
	Secret string `json:"client_secret"`
}

func getClient() client {
	return client{
		ID:     os.Getenv("BGGH_CLIENT"),
		Secret: os.Getenv("BGGH_SECRET"),
	}
}
