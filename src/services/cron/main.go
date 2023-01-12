package main

import (
	"ecomdream/src/services/cron/jobs/dangling_prompts"
	"ecomdream/src/services/cron/jobs/push_versions"
	"time"
)

type Job interface {
	Start(interval time.Duration)
}

func GetJob(jobName string) Job {
	switch jobName {
		case "push_versions":
			return push_versions.CreatJob()
		case "dangling_prompts":
			return dangling_prompts.CreatJob()
		default:
			return nil
	}
}

func main() {
	go GetJob("push_versions").Start(1 * time.Minute)
	go GetJob("dangling_prompts").Start(30 * time.Second)

	<- make(chan uint8)
}
