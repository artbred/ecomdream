package redisdb

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type Key string

func (k Key) IsBlocked() bool {
	rdb := Connection()

	isBlocked, err := rdb.Get(context.Background(), string(k)).Bool()
	if err != nil {
		if err == redis.Nil {
			return false
		} else {
			logrus.Error(err)
			return false
		}
	}

	return isBlocked
}

func BuildFreezeEndpointKey(endpoint, id string) Key {
	return Key(fmt.Sprintf("freeze_api_endpoint:%s:%s", endpoint, id))
}

func BuildBlockReplicatePrediction(predictionID string) Key {
	return Key(fmt.Sprintf("block_replicate_prediction:%s", predictionID))
}
