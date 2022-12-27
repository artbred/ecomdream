package dangling_prompts

import (
	"ecomdream/src/services/worker/logger"
	"github.com/sirupsen/logrus"
	"time"
)

const jobName = "dangling_prompts"

type PromptsJob struct {
	JobName string
	logger  *logrus.Entry
}

func (j *PromptsJob) Start(interval time.Duration) {
	for {
		j.Logic()
		time.Sleep(interval)
	}
}

func CreatJob() *PromptsJob {
	return &PromptsJob{
		JobName: jobName,
		logger:  logger.CreateLogger(jobName),
	}
}

