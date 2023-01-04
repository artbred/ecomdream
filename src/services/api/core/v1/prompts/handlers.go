package prompts

import (
	"context"
	"ecomdream/src/domain/models"
	"ecomdream/src/domain/replicate"
	"ecomdream/src/pkg/storages/redisdb"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/twinj/uuid"
	"time"
)

// CreatePromptForIDHandler handler that start prediction for model
// @Description Start prediction for prompt
// @Summary Start prediction for prompt
// @Tags prompts
// @Accept json
// @Produce json
// @Param id path string true "Version ID"
// @Param prompt_data body CreatePromptRequest true "Prompt data"
// @Success 201 {object} CreatePromptRequest
// @Router /v1/prompts/create/{id} [post]
func CreatePromptForIDHandler(ctx *fiber.Ctx) error {
	req := &CreatePromptRequest{}

	if err := req.Validate(ctx); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": err.Error(),
		})
	}

	version, err := models.GetVersion(req.VersionID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Please try again later",
		})
	}

	if version == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid version id",
		})
	}

	if version.PushedAt == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Please wait for your model to be ready",
		})
	}

	if version.DeletedAt != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":    fiber.StatusForbidden,
			"message": "Your model has been deleted",
		})
	}

	//hasRunningPrompts, err := version.HasRunningPrompts(); if err != nil {
	//	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	//		"code":    fiber.StatusInternalServerError,
	//		"message": "Please try again later",
	//	})
	//}
	//
	//if hasRunningPrompts {
	//	return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
	//		"code":    fiber.StatusForbidden,
	//		"message": "You have one prompt running, the max wait time is 5 minutes, please wait for it to complete",
	//	})
	//}

	features, err := version.GetUnifiedFeatures(); if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Please try again later",
		})
	}

	if features.FeatureAmountImages - version.AmountImagesGenerated - req.AmountImages <= 0 {
		text := fmt.Sprintf("You can only generate %d images", features.FeatureAmountImages - version.AmountImagesGenerated)
		if features.FeatureAmountImages - version.AmountImagesGenerated <= 0 {
			text = "You ran out of images, you need to purchase extra package to continue working"
		}
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":    fiber.StatusForbidden,
			"message": text,
		})
	}

	prompt := &models.Prompt{
		ID:             uuid.NewV4().String(),
		VersionID:      version.ID,
		PromptText:     req.Prompt,
		PromptNegative: req.NegativePrompt,
		AmountImages:   req.AmountImages,
		InferenceSteps: 50,
		Width:          512,
		Height:         512,
		PromptStrength: 0.8,
		GuidanceScale:  7.5,
	}

	replicateReq, err := prompt.TransferToReplicateBody(version); if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Please try again later",
		})
	}

	replicateInitResponse, err := replicate.StartPrediction(replicateReq)
	if err != nil {
		logrus.WithError(err).Errorf("Can't start prediction for prompt %+v", prompt)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Please try again later",
		})
	}

	prompt.PredictionID = replicateInitResponse.ID
	err = prompt.Create(); if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Please try again later",
		})
	}

	key := redisdb.BuildBlockReplicatePrediction(prompt.PredictionID)
	rdb := redisdb.Connection()
	rdb.SetNX(context.Background(), string(key), true, 5*time.Minute)
	defer rdb.Del(context.Background(), string(key))

	replicateOutResponse, err := replicate.WaitForPrediction(context.Background(), prompt.PredictionID)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to create images for prompt %s", prompt.ID)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Please try again later",
		})
	}

	imagesGeneratedUrls, err := ReplicateToCloudflare(replicateOutResponse, prompt)
	if err != nil || len(imagesGeneratedUrls) == 0 {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Please try again later",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(CreatePromptResponse{
		Code:           fiber.StatusCreated,
		Images:         imagesGeneratedUrls,
		ImagesLeft:     features.FeatureAmountImages - version.AmountImagesGenerated - len(imagesGeneratedUrls),
		PromptText:     prompt.PromptText,
		PromptNegative: prompt.PromptNegative,
	})
}


// ListPromptsForIDHandler handler that returns prompts and images for version
// @Description Returns prompts and images for version
// @Summary Returns prompts and images for version
// @Tags prompts
// @Accept json
// @Produce json
// @Param id path string true "Version ID"
// @Success 201 {object} CreatePromptRequest
// @Router /v1/prompts/list/{id} [get]
func ListPromptsForIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id"); if len(id) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": fiber.StatusBadRequest,
			"message": "Provide version id",
		})
	}

	prompts, err := models.GetCompletedPromptsForVersion(id); if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": fiber.StatusInternalServerError,
			"message": "Please try again later",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(ListPromptsResponse{
		Code: fiber.StatusOK,
		Prompts: prompts,
	})
}
