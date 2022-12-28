package dangling_prompts

import (
	"context"
	"ecomdream/src/domain/models"
	"ecomdream/src/domain/replicate"
	"ecomdream/src/pkg/storages/redisdb"
	"ecomdream/src/services/api/core/v1/prompts"
)

func (j *PromptsJob) Logic() {
	runningPrompts, err := models.GetRunningPrompts(); if err != nil {
		return
	}

	counterUnblocked := 0
	counterPushed := 0

	for _, prompt := range runningPrompts {
		key := redisdb.BuildBlockReplicatePrediction(prompt.PredictionID)
		if key.IsBlocked() {
			continue
		}

		counterUnblocked++

		result, err := replicate.CheckPrediction(context.Background(), prompt.PredictionID)
		if err != nil {
			j.logger.WithError(err).Errorf("can't check prediction status for prompt %s", prompt.ID)
			continue
		}

		if result == nil {
			continue
		}

		_, err = prompts.ReplicateToCloudflare(result, prompt)
		if err == nil  {
			counterPushed++
		}
	}

	j.logger.Infof("Pushed %d/%d prompts", counterPushed, counterUnblocked)
}
