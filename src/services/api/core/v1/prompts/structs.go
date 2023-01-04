package prompts

import (
	"ecomdream/src/domain/models"
	"errors"
	"github.com/gofiber/fiber/v2"
)

type CreatePromptRequest struct {
	VersionID string `json:"-"`
	Prompt string `json:"prompt"`
	NegativePrompt *string `json:"negative_prompt"`
	AmountImages int `json:"amount_images"`
}

type CreatePromptResponse struct {
	Code int `json:"code"`
	Images []string `json:"images"`
	ImagesLeft int `json:"images_left"`
	PromptText string `json:"prompt_text"`
	PromptNegative *string `json:"prompt_negative"`
}

type ListPromptsResponse struct {
	Code int `json:"code"`
	Prompts []models.Prompt `json:"prompts"`
}

func (r *CreatePromptRequest) Validate(ctx *fiber.Ctx) (err error) {
	err = ctx.BodyParser(r); if err != nil {
		return
	}

	versionID := ctx.Params("id"); if len(versionID) == 0 {
		return errors.New("Provide version id")
	}

	r.VersionID = versionID

	if len(r.Prompt) == 0 {
		return errors.New("Provide prompt")
	} else if len(r.Prompt) > 165 {
		return errors.New("Maximum prompt length is 165")
	}

	if r.NegativePrompt != nil {
		if len(*r.NegativePrompt) == 0 {
			r.NegativePrompt = nil
		}
	}

	if r.AmountImages != 1 && r.AmountImages != 4 {
		return errors.New("You can generate only 1 or 4 images")
	}

	return nil
}
