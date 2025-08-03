package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()

	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		PrettyPrint:     true,
	})

	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.DebugLevel)
}

func GetLogger() *logrus.Logger {
	return Logger
}
