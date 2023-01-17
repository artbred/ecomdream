package replicate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

var (
	dreamBoothBaseURL = "https://dreambooth-api-experimental.replicate.com/v1/trainings"
	TrainerVersion = "d5e058608f43886b9620a8fbb1501853b8cbae4f45c857a014011c86ee614ffb"
)

type DreamBoothRequest struct {
	Input DreamBoothInput `json:"input"`
	Model string          `json:"model"`
}

type DreamBoothInput struct {
	InstancePrompt string `json:"instance_prompt"`
	ClassPrompt string `json:"class_prompt"`
	InstanceData string `json:"instance_data"`
	MaxTrainSteps int64 `json:"max_train_steps"`
	TrainerVersion string `json:"trainer_version"`
}

type DreamBoothResponse struct {
	ID    string          `json:"id"`
	Input DreamBoothInput `json:"input"`
	Model string          `json:"model"`
	Status string `json:"status"`
	Version *string `json:"version"`
}

var (
	MaxTrainingSteps int64 = 2000
	ModelName = "artbred/ecomdream"
	UniqueIdentifier = "xjy"
)

func ConstructDreamBoothInputs(class, zipURL, trainerVersion string) DreamBoothRequest {
	return DreamBoothRequest{
		Input: DreamBoothInput{
			InstanceData:   zipURL,
			MaxTrainSteps:  MaxTrainingSteps,
			ClassPrompt:    fmt.Sprintf("a %s", class),
			InstancePrompt: fmt.Sprintf("a photo of a %s %s", UniqueIdentifier, class),
			TrainerVersion: trainerVersion,
		},
		Model: ModelName,
	}
}


func CheckDreamBoothTraining(ctx context.Context, id string) (response DreamBoothResponse, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/%s", dreamBoothBaseURL, id), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", apiToken))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req); if err != nil {
		return
	}

	defer res.Body.Close()

	resBytes, err := io.ReadAll(res.Body); if err != nil {
		return
	}

	statusOK := res.StatusCode >= 200 && res.StatusCode < 300
	if !statusOK {
		err = fmt.Errorf("replicate status code is not 2xx, %s", string(resBytes))
		return
	}

	err = json.Unmarshal(resBytes, &response); if err != nil {
		return
	}

	return
}

func StartDreamBoothTraining(ctx context.Context, modelRequest DreamBoothRequest) (response DreamBoothResponse, err error) {
	body, err := json.Marshal(modelRequest); if err != nil {
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", dreamBoothBaseURL, bytes.NewReader(body))
	if err != nil {
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("Token %s", apiToken))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req); if err != nil {
		return
	}

	defer res.Body.Close()

	resBytes, err := io.ReadAll(res.Body); if err != nil {
		return
	}

	statusOK := res.StatusCode >= 200 && res.StatusCode < 300
	if !statusOK {
		err = fmt.Errorf("replicate status code is not 2xx, %s", string(resBytes))
		return
	}

	err = json.Unmarshal(resBytes, &response)
	return
}

func init() {
	MaxTrainingSteps, _ = strconv.ParseInt(os.Getenv("MAX_TRAINING_STEPS"), 10, 64)
}
