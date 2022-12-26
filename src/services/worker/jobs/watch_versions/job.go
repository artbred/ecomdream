package watch_versions

import (
	"ecomdream/src/services/worker/logger"
	"github.com/sirupsen/logrus"
	"time"
)

const jobName = "watch-versions"

type VersionsJob struct {
	JobName string
	logger  *logrus.Entry
}

func (j *VersionsJob) Start(interval time.Duration) {
	for {
		j.Logic()
		time.Sleep(interval)
	}
}

func CreatJob() *VersionsJob {
	return &VersionsJob{
		JobName: jobName,
		logger:  logger.CreateLogger(jobName),
	}
}
