package log

import (
	"github.com/sirupsen/logrus"
)

var Logger = GetLogger("info")

func GetLogger(logLv string) *logrus.Logger {
	logger := logrus.New()

	textFormatter := new(logrus.TextFormatter)
	textFormatter.TimestampFormat = "2006-01-02 15:04:05"
	textFormatter.FullTimestamp = true
	textFormatter.ForceColors = true

	logger.SetFormatter(textFormatter)

	logLevel, err := logrus.ParseLevel(logLv)
	if err != nil {
		logger.SetLevel(logrus.InfoLevel)
	}
	logger.SetLevel(logLevel)

	return logger
}
