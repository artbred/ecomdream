package main

import (
	"ecomdream/src/services/worker/jobs/watch_versions"
	"time"
)

type Job interface {
	Start(interval time.Duration)
}

func GetJob(jobName string) Job {
	switch jobName {
	case "watch_versions":
		return watch_versions.CreatJob()
	default:
		return nil
	}
}

func main() {
	go GetJob("watch_versions").Start(1 * time.Minute)

	<-make(chan bool)
}
