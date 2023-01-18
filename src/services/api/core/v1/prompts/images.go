package prompts

import (
	"ecomdream/src/domain/models"
	"ecomdream/src/domain/replicate"
	"ecomdream/src/pkg/config"
	"ecomdream/src/pkg/storages/cloudflare"
	"github.com/sirupsen/logrus"
	"sync"
)

func TransferReplicateImagesToCloudflareAndSave(replicateOutResponse *replicate.Response, prompt *models.Prompt) ([]string, error) {
	var imagesGeneratedUrls []string
	var wg sync.WaitGroup

	imageChan := make(chan *models.Image, len(replicateOutResponse.Output))

	for _, imageReplicate := range replicateOutResponse.Output {
		wg.Add(1)
		go func(prompt *models.Prompt, imageReplicateUrl string) {
			defer wg.Done()
			image := replicateImageToCloudflare(prompt, imageReplicateUrl)
			if image != nil {
				imageChan <- image
			}
		}(prompt, imageReplicate)
	}

	wg.Wait()
	close(imageChan)

	for image := range imageChan {
		imagesGeneratedUrls = append(imagesGeneratedUrls, image.CdnURL)
	}

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
