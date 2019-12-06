package secrets

import (
	"context"
	"errors"

	"github.com/Squwid/bytegolf/firestore"
	"github.com/mitchellh/mapstructure"
)

const collection = "secrets"

// Client is a structure that is stored in cloud for api secrets rather than
// passing environmental variables around
type Client struct {
	Client string
	Secret string
}

// Store ...
func (c *Client) Store(key string) error {
	ctx := context.Background()
	if c == nil {
		return errors.New("Error client was null")
	}
	_, err := firestore.Client.Collection(collection).Doc(key).Set(ctx, *c)
	return err
}

// GetClient ...
func GetClient(key string) (*Client, error) {
	ctx := context.Background()
	ref, err := firestore.Client.Collection(collection).Doc(key).Get(ctx)
	if err != nil {
		return nil, err
	}
	var c Client
	err = mapstructure.Decode(ref.Data(), &c)
	if err != nil {
		return nil, err
	}
	return &c, err
}

// Must is a wrapper around Get to panic if there is an error for easy entrance to
// init functions
func Must(s interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return s
}
