package tests

import (
	"context"
	"testing"

	"github.com/Squwid/bytegolf/lib/api"
	"github.com/Squwid/bytegolf/lib/sqldb"
	"github.com/sirupsen/logrus"
)

func TestPopulateLanguages(t *testing.T) {
	// Simple map to hold most common languages for testing.

	var langs = []api.LanguageDB{
		{
			Language:  "Go",
			Version:   "1.19.4",
			Image:     "golang:1.19.4-alpine3.17",
			Active:    true,
			Cmd:       "go run",
			Extension: ".go",
		},
		{
			Language:  "Python",
			Version:   "3.11.1",
			Image:     "python:3.11.1-alpine3.17",
			Active:    true,
			Cmd:       "python",
			Extension: ".py",
		},
	}

	// Loop through the map and insert each language into the database.
	for _, lang := range langs {
		if _, err := sqldb.DB.NewInsert().Model(&lang).
			Exec(context.Background()); err != nil {
			logrus.WithError(err).Errorf("Error writing language %v\n",
				lang.Language)
		}
	}
}
