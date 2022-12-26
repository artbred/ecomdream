package versions

import (
	"ecomdream/src/pkg/api/replicate"
	"fmt"
	"os"
	"strconv"
)

var (
	maxTrainingSteps int64 = 2000
	modelName        string
	uniqueIdentifier = "xjy"
)

func ConstructDreamBoothInputs(class, zipURL string) replicate.DreamBoothRequest {
	return replicate.DreamBoothRequest{
		Input: replicate.DreamBoothInput{
			InstanceData:   zipURL,
			MaxTrainSteps:  maxTrainingSteps,
			ClassPrompt:    fmt.Sprintf("a %s", class),
			InstancePrompt: fmt.Sprintf("a photo of a %s %s", uniqueIdentifier, class),
		},
		Model: modelName,
	}
}

func init() {
	maxTrainingSteps, _ = strconv.ParseInt(os.Getenv("MAX_TRAINING_STEPS"), 10, 64)
	modelName = os.Getenv("REPLICATE_MODEL_NAME")
}
