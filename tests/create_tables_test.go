package tests

import (
	"context"
	"testing"

	"github.com/Squwid/bytegolf/lib/api"
	"github.com/Squwid/bytegolf/lib/auth"
	"github.com/Squwid/bytegolf/lib/sqldb"
	"github.com/sirupsen/logrus"
)

// TestCreateTables will create all sql tables based on Go models.
// https://bun.uptrace.dev/guide/models.html
func TestCreateTables(t *testing.T) {
	// t.SkipNow()
	ctx := context.Background()

	// Create 'users' table.
	if _, err := sqldb.DB.NewCreateTable().
		Model((*auth.BytegolfUserDB)(nil)).Exec(ctx); err != nil {
		logrus.WithError(err).Warnf("Skipping creation of 'users' table.")
	} else {
		logrus.Infof("Successfully created 'users' table.")
	}

	// Create 'holes' table.
	if _, err := sqldb.DB.NewCreateTable().
		Model((*api.HoleDB)(nil)).Exec(ctx); err != nil {
		logrus.WithError(err).Warnf("Skipping creation of 'holes' table.")
	} else {
		logrus.Infof("Successfully created 'holes' table.")
	}

	// Create 'tests' table.
	if _, err := sqldb.DB.NewCreateTable().
		Model((*api.TestDB)(nil)).Exec(ctx); err != nil {
		logrus.WithError(err).Warnf("Skipping creation of 'tests' table.")
	} else {
		logrus.Infof("Successfully created 'tests' table.")
	}

	// Create 'languages' table.
	if _, err := sqldb.DB.NewCreateTable().
		Model((*api.LanguageDB)(nil)).Exec(ctx); err != nil {
		logrus.WithError(err).Warnf("Skipping creation of 'languages' table.")
	} else {
		logrus.Infof("Successfully created 'languages' table.")
	}

	// Create 'submissions' table.
	if _, err := sqldb.DB.NewCreateTable().
		Model((*api.SubmissionDB)(nil)).Exec(ctx); err != nil {
		logrus.WithError(err).Warnf("Skipping creation of 'submissions' table.")
	} else {
		logrus.Infof("Successfully created 'submissions' table.")
	}

	// Create 'jobs' table.
	if _, err := sqldb.DB.NewCreateTable().
		Model((*api.JobOutputDB)(nil)).Exec(ctx); err != nil {
		logrus.WithError(err).Warnf("Skipping creation of 'jobs' table.")
	} else {
		logrus.Infof("Successfully created 'jobs' table.")
	}

}
