package logger

import (
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type CustomJSONFormatter struct {
	logrus.JSONFormatter
}

func (f *CustomJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Ensure the entry's time is in UTC
	entry.Time = entry.Time.UTC()

	originalFormatted, err := f.JSONFormatter.Format(entry)
	if err != nil {
		return nil, err
	}

	formattedString := string(originalFormatted)
	if strings.Contains(formattedString, "+00:00") {
		formattedString = strings.Replace(formattedString, "+00:00", "Z", 1)
	} else {
		formattedString = strings.Replace(formattedString, "-04:00", "Z", 1)
	}

	return []byte(formattedString), nil
}

var Logger = logrus.New()

func init() {
	Logger.Formatter = &CustomJSONFormatter{
		JSONFormatter: logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		},
	}
	Logger.SetLevel(logrus.DebugLevel)

	file, err := os.OpenFile("/var/log/webapp/webapp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Logger.Out = file
	} else {
		Logger.Info("Failed to log to file, using default stderr")
	}
}
