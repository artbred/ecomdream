package prompts

import (
	"ecomdream/src/domain/models"
	"github.com/gofiber/fiber/v2"
)

// CreatePromptHandler handler that start prediction for model
// @Description Start prediction for prompt
// @Summary Start prediction for prompt
// @Tags prompts
// @Accept json
// @Produce json
// @Param id query string true "Version ID"
// @Param prompt_data body CreatePromptRequest true "Prompt data"
// @Success 201 {object} CreatePromptRequest
// @Router /v1/prompts/create [post]
func CreatePromptHandler(ctx *fiber.Ctx) error {
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

	features, err := version.GetUnifiedFeatures()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Please try again later",
		})
	}

	_ = features

	return nil
}
