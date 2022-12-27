package dangling_prompts

import (
	"context"
	"ecomdream/src/domain/models"
	"ecomdream/src/domain/replicate"
	"ecomdream/src/pkg/storages/redisdb"
	"ecomdream/src/services/api/core/v1/prompts"
	"github.com/go-redis/redis/v8"
)

func (j *PromptsJob) Logic() {
	runningPrompts, err := models.GetRunningPrompts(); if err != nil {
		return
	}

	counterUnblocked := 0
	counterPushed := 0

	for _, prompt := range runningPrompts {
		key := redisdb.BuildReplicatePredictionFreeze(prompt.PredictionID)
		rdb := redisdb.Connection()

		isBlocked, err := rdb.Get(context.Background(), key).Bool()
		if err != nil {
			if err == redis.Nil {
				isBlocked = false
			} else {
				j.logger.Error(err)
			}
		}

		if isBlocked {
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
