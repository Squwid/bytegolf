package tests

import (
	"context"
	"testing"

	"github.com/Squwid/bytegolf/auth"
	"github.com/Squwid/bytegolf/sqldb"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := sqldb.Open(); err != nil {
		logrus.WithError(err).Fatalf("Error connecting to db")
	}
}

// TestCreateTables will create all sql tables based on Go models.
// https://bun.uptrace.dev/guide/models.html
func TestCreateTables(t *testing.T) {
	// t.SkipNow()

	ctx := context.Background()
	_, err := sqldb.DB.NewCreateTable().Model((*auth.BytegolfUserDB)(nil)).Exec(ctx)
	if err != nil {
		logrus.WithError(err).Fatalf("Error creating user table")
	}
}
