package docker

import (
	"context"
	"encoding/json"
	"io"

	"github.com/Squwid/bytegolf/lib/api"
	"github.com/Squwid/bytegolf/lib/log"
	"github.com/Squwid/bytegolf/lib/sqldb"
	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
)

type jsonMessage struct {
	Status string `json:"status,omitempty"`
}

// SyncImages syncs the images in the database with the images on the docker host.
func SyncImages(ctx context.Context) error {
	logger := log.GetLogger().WithField("Action", "SyncImages")
	logger.Infof("Syncing images. This may take a second...")

	var langs []api.LanguageDB
	if err := sqldb.DB.NewSelect().Model(&langs).Scan(ctx); err != nil {
		return errors.Wrap(err, "Error retrieving languages")
	}

	var images []string

	for _, lang := range langs {
		reader, err := Client.c.ImagePull(ctx, lang.Image, types.ImagePullOptions{})
		if err != nil {
			log.GetLogger().WithField("Action", "SyncImages").Errorf("Error pulling image %s")
			continue
		}

		decoder := json.NewDecoder(reader)
		for {
			var msg jsonMessage
			if err := decoder.Decode(&msg); err != nil {
				if err == io.EOF {
					break
				}
				logger.WithError(err).Errorf("Error reading image pull output for %s", lang.Image)
				break
			}

			logger.Debugf("%s", msg.Status)
		}
		reader.Close()
		images = append(images, lang.Image)
	}
	logger.WithField("ImagesSynced", len(images)).Infof("Done syncing images")

	return nil
}
