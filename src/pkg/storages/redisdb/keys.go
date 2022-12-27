package redisdb

import "fmt"

//TODO key interface is blocked

func BuildFreezeEndpointKey(endpoint, id string) string {
	return fmt.Sprintf("freeze_api_endpoint:%s:%s", endpoint, id)
}

func BuildReplicatePredictionFreeze(predictionID string) string {
	return fmt.Sprintf("freeze_replicate_prediction:%s", predictionID)
}
