package replicate

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func StartPrediction(body []byte) (response *Response, err error) {
	req, err := http.NewRequest("POST", baseURL, bytes.NewReader(body))
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

	bodyResponse, err := io.ReadAll(res.Body); if err != nil {
		return
	}

	if res.StatusCode != 201 {
		err = errors.New(fmt.Sprintf("replicate status code is not 201, %s", string(bodyResponse)))
		return
	}

	err = json.Unmarshal(bodyResponse, &response)
	return
}

func CheckPrediction(ctx context.Context, predictionID string) (result *Response, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/%s", baseURL, predictionID), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", apiToken))
	req.Header.Add("Content-Type", "application/json")

	httpClient := http.Client{}

	res, err := httpClient.Do(req); if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body); if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Status == "succeeded" {
		return result, nil
	}

	if result.Status == "failed" {
		return nil, fmt.Errorf("failed to inference, %s", result.Error)
	}

	return nil, nil
}

func WaitForPrediction(ctx context.Context, predictionID string) (result *Response, err error) {
	errorChan := make(chan error)

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/%s", baseURL, predictionID), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", apiToken))
	req.Header.Add("Content-Type", "application/json")

	httpClient := http.Client{}

	go func() {
		for {
			res, err := httpClient.Do(req); if err != nil {
				errorChan <- err
				break
			}

			body, err := io.ReadAll(res.Body); if err != nil {
				res.Body.Close()
				errorChan <- err
				break
			}

			res.Body.Close()

			if err = json.Unmarshal(body, &result); err != nil {
				errorChan <- err
				break
			}

			if result.Status == "succeeded" {
				errorChan <- nil
				break
			}

			if result.Status == "failed" {
				errorChan <- errors.New(fmt.Sprintf("failed to inference, %s", result.Error))
				break
			}
		}
	}()

	err = <-errorChan; if err != nil {
		return nil, parseError(err)
	}

	return
}
