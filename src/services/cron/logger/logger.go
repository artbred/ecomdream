package logger

import (
	"github.com/sirupsen/logrus"
)

func CreateLogger(jobName string) *logrus.Entry {
	logger := logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return logger.WithField("job", jobName)
}
