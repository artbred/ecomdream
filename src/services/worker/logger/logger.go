package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
	"strconv"
)

func CreateLogger(jobName string) *logrus.Entry {
	logger := logrus.New()

	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
			return "", fmt.Sprintf("[%s]", fileName)
		},
		FullTimestamp: true,
	})

	return logger.WithField("job", jobName)
}
