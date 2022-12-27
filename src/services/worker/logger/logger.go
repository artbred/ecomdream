package logger

import (
	"github.com/sirupsen/logrus"
)

func CreateLogger(jobName string) *logrus.Entry {
	logger := logrus.New()
	return logger.WithField("job", jobName)
}
