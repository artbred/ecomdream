package prompts

import (
	"ecomdream/src/domain/models"
	"ecomdream/src/domain/replicate"
	"ecomdream/src/pkg/config"
	"ecomdream/src/pkg/storages/cloudflare"
	"github.com/sirupsen/logrus"
	"sync"
)


func ReplicateImageToCloudflare(replicateOutResponse *replicate.Response, prompt *models.Prompt) ([]string, error) {
	var imagesGeneratedUrls []string
	var wg sync.WaitGroup

	imageChan := make(chan *models.Image)

	for _, imageReplicate := range replicateOutResponse.Output {
		wg.Add(1)
		go func(prompt *models.Prompt, imageReplicate string) {
			defer wg.Done()
			imageChan <- replicateImageToCloudflare(prompt, imageReplicate)
		}(prompt, imageReplicate)
	}

	go func() {
		for i := 0; i < len(replicateOutResponse.Output); i++ {
			select {
			case image := <-imageChan:
				if image != nil {
					imagesGeneratedUrls = append(imagesGeneratedUrls, image.CdnURL)
				}
			}
		}
	}()

	wg.Wait()

	prompt.PredictionTime = &replicateOutResponse.Metrics.PredictTime
	return imagesGeneratedUrls, prompt.MarkAsFinished()
}

func replicateImageToCloudflare(prompt *models.Prompt, urlReplicate string) *models.Image {
	imageCloudflare, err := cloudflare.UploadImageByURL(cloudflare.ImageUploadRequestByURL{
		URL: urlReplicate,
		RequireSignedURLs: false,
		Metadata: map[string]interface{}{
			"version_id": prompt.VersionID,
			"debug": config.Debug,
		},
	})

	if err != nil {
		logrus.WithError(err).Errorf("can't transfer image from replicate to cloudflare")
		return nil
	}

	image := &models.Image{
		ID: imageCloudflare.Result.ID,
		PromptID: prompt.ID,
		CdnURL: imageCloudflare.Result.Variants[0],
		Width: prompt.Width,
		Height: prompt.Height,
	}

	err = image.Create(); if err != nil {
		return nil
	}

	return image
}
