package tests

import (
	"context"
	"testing"

	"github.com/Squwid/bytegolf/api"
	"github.com/Squwid/bytegolf/auth"
	"github.com/Squwid/bytegolf/sqldb"
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
		logrus.WithError(err).Warnf("Skipping creation of 'users' table")
	}

	// Create 'holes' table
	if _, err := sqldb.DB.NewCreateTable().
		Model((*api.HoleDB)(nil)).Exec(ctx); err != nil {
		logrus.WithError(err).Warnf("Skipping creation of 'holes' table")
	}

}
