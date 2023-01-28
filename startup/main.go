package main

import (
	"context"
	"io"
	"os"

	"github.com/Squwid/bytegolf/lib/api"
	"github.com/Squwid/bytegolf/lib/auth"
	"github.com/Squwid/bytegolf/lib/sqldb"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func init() {
	if err := sqldb.Open(); err != nil {
		logrus.WithError(err).Fatalf("Error connecting to db")
	}
}

func main() {
	createTables()
	populateLanguages()
	populateHoles()
	populateTests()
}

func createTables() {
	ctx := context.Background()

	// Create 'users' table.
	if _, err := sqldb.DB.NewCreateTable().
		Model((*auth.BytegolfUserDB)(nil)).Exec(ctx); err != nil {
		// logrus.WithError(err).Warnf("Skipping creation of 'users' table.")
	} else {
		logrus.Infof("Successfully created 'users' table.")
	}

	// Create 'holes' table.
	if _, err := sqldb.DB.NewCreateTable().
		Model((*api.HoleDB)(nil)).Exec(ctx); err != nil {
		// logrus.WithError(err).Warnf("Skipping creation of 'holes' table.")
	} else {
		logrus.Infof("Successfully created 'holes' table.")
	}

	// Create 'tests' table.
	if _, err := sqldb.DB.NewCreateTable().
		Model((*api.TestDB)(nil)).Exec(ctx); err != nil {
		// logrus.WithError(err).Warnf("Skipping creation of 'tests' table.")
	} else {
		logrus.Infof("Successfully created 'tests' table.")
	}

	// Create 'languages' table.
	if _, err := sqldb.DB.NewCreateTable().
		Model((*api.LanguageDB)(nil)).Exec(ctx); err != nil {
		// logrus.WithError(err).Warnf("Skipping creation of 'languages' table.")
	} else {
		logrus.Infof("Successfully created 'languages' table.")
	}

	// Create 'submissions' table.
	if _, err := sqldb.DB.NewCreateTable().
		Model((*api.SubmissionDB)(nil)).Exec(ctx); err != nil {
		// logrus.WithError(err).Warnf("Skipping creation of 'submissions' table.")
	} else {
		logrus.Infof("Successfully created 'submissions' table.")
	}

	// Create 'jobs' table.
	if _, err := sqldb.DB.NewCreateTable().
		Model((*api.JobOutputDB)(nil)).Exec(ctx); err != nil {
		// logrus.WithError(err).Warnf("Skipping creation of 'jobs' table.")
	} else {
		logrus.Infof("Successfully created 'jobs' table.")
	}
}

func populateLanguages() {
	file, err := os.Open("languages.yaml")
	if err != nil {
		logrus.WithError(err).Fatalf("Error opening languages.yaml")
	}
	defer file.Close()

	bs, err := io.ReadAll(file)
	if err != nil {
		logrus.WithError(err).Fatalf("Error reading languages.yaml")
	}

	var languages []api.LanguageDB
	if err := yaml.Unmarshal(bs, &languages); err != nil {
		logrus.WithError(err).Fatalf("Error unmarshaling languages.yaml")
	}

	for _, lang := range languages {
		if _, err := sqldb.DB.NewInsert().
			Model(&lang).
			On("CONFLICT (id) DO UPDATE").
			Exec(context.Background()); err != nil {
			logrus.WithError(err).Fatalf("Error inserting language %s", lang.Language)
		}
	}

	logrus.Infof("Successfully inserted %v languages.", len(languages))
}
