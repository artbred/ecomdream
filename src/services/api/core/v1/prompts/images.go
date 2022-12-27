package prompts

import (
	"ecomdream/src/domain/models"
	"ecomdream/src/pkg/configs"
	"ecomdream/src/pkg/storages/cloudflare"
	"github.com/sirupsen/logrus"
)

func replicateImageToCloudflare(prompt *models.Prompt, urlReplicate string, imageChan chan *models.Image)  {
	imageCloudflare, err := cloudflare.UploadImageByURL(cloudflare.ImageUploadRequestByURL{
		URL: urlReplicate,
		RequireSignedURLs: false,
		Metadata: map[string]interface{}{
			"version_id": prompt.VersionID,
			"debug": configs.Debug,
		},
	})

	if err != nil {
		logrus.WithError(err).Errorf("can't transfer image from replicate to cloudflare")
		imageChan <- nil
		return
	}

	image := &models.Image{
		ID: imageCloudflare.Result.ID,
		PromptID: prompt.ID,
		CdnURL: imageCloudflare.Result.Variants[0],
		Width: prompt.Width,
		Height: prompt.Height,
	}

	err = image.Create(); if err != nil {
		imageChan <- nil
		return
	}

	imageChan <- image
}
