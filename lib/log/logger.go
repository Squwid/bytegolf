package log

import (
	"github.com/Squwid/bytegolf/lib/globals"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	// logger.SetFormatter(&logrus.JSONFormatter{})

	if globals.IsDebugMode() {
		logger.SetLevel(logrus.DebugLevel)
	}
}

// GetLogger returns the logger instance for logging.
func GetLogger() *logrus.Logger {
	return logger
}

// NewEntry returns a new logrus entry
func NewEntry() *logrus.Entry {
	return logrus.NewEntry(GetLogger())
}
