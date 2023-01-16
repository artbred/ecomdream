package dangling_prompts

import (
	"context"
	"ecomdream/src/domain/models"
	"ecomdream/src/domain/replicate"
	"ecomdream/src/pkg/storages/redisdb"
	"ecomdream/src/services/api/core/v1/prompts"
	"github.com/sirupsen/logrus"
)

func (j *PromptsJob) Logic() {
	runningPrompts, err := models.GetRunningPrompts(); if err != nil {
		return
	}

	counterRunningPrompts := len(runningPrompts)
	counterAccessiblePrompts, counterPushedPrompts := 0, 0

	for _, prompt := range runningPrompts {
		key := redisdb.BuildFreezeReplicatePrediction(prompt.PredictionID)
		if key.IsFrozen() {
			continue
		}

		counterAccessiblePrompts++

		result, err := replicate.CheckPrediction(context.Background(), prompt.PredictionID)
		if err != nil {
			j.logger.WithError(err).Errorf("can't check prediction status for prompt %s", prompt.ID)
			continue
		}

		if result == nil {
			continue
		}

		_, err = prompts.TransferReplicateImagesToCloudflareAndSave(result, prompt)
		if err == nil  {
			counterPushedPrompts++
		} else {
			logrus.WithError(err).Error("can't transfer ")
		}
	}

	j.logger.Infof("Pushed %d/%d prompts, from %d prompts", counterAccessiblePrompts, counterAccessiblePrompts, counterRunningPrompts)
}
