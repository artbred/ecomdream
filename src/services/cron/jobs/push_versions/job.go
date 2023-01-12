package push_versions

import (
	"ecomdream/src/services/cron/logger"
	"github.com/sirupsen/logrus"
	"time"
)

const jobName = "push_versions"

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
