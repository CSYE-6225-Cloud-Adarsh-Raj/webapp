package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func init() {
	Logger.Formatter = &logrus.JSONFormatter{}

	file, err := os.OpenFile("/var/log/webapp/webapp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Logger.Out = file
	} else {
		Logger.Info("Failed to log to file, using default stderr")
	}
}
