package tests

import (
	"github.com/Squwid/bytegolf/lib/sqldb"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := sqldb.Open(); err != nil {
		logrus.WithError(err).Fatalf("Error connecting to db")
	}
}
