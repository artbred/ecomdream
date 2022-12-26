package replicate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	dreamBoothBaseURL = "https://dreambooth-api-experimental.replicate.com/v1/trainings"
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
}

type DreamBoothResponse struct {
	ID    string          `json:"id"`
	Input DreamBoothInput `json:"input"`
	Model string          `json:"model"`
	Status string `json:"status"`
	Version *string `json:"version"`
}

func CheckDreamBoothTraining(ctx context.Context, id string) (response DreamBoothResponse, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/%s", dreamBoothBaseURL, id), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", apiToken))

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

	err = json.Unmarshal(resBytes, &response); if err != nil {
		return
	}

	return
}

