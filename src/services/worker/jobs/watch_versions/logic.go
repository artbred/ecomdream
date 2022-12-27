package watch_versions

import (
	"context"
	"ecomdream/src/domain/models"
	"ecomdream/src/domain/replicate"
	"time"
)

func (j *VersionsJob) Logic() {
	runningVersions, err := models.GetRunningVersions()
	if err != nil {
		return
	}

	countFinished := 0

	for _, version := range runningVersions {
		res, err := replicate.CheckDreamBoothTraining(context.Background(), version.PredictionID)
		if err != nil {
			j.logger.WithField("version_id", version.ID).Error(err)
			continue
		}

		if res.Status == "succeeded" && res.Version != nil {
			err := version.MarkAsPushed(res.Version)
			if err == nil {
				countFinished++
			}
		} else if time.Now().UTC().Sub(version.CreatedAt) >= 1*time.Hour {
			j.logger.WithField("version_id", version.ID).Warning("Takes too long, need to check")
			// TODO informer
		}
	}

	j.logger.Infof("Marked %d/%d versions as pushed", countFinished, len(runningVersions))
}
