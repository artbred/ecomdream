package prompts

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

type CreatePromptRequest struct {
	VersionID string `json:"version_id"`
	Prompt string `json:"prompt"`
	AmountImages int `json:"amount_images"`
}

func (r *CreatePromptRequest) Validate(ctx *fiber.Ctx) (err error) {
	err = ctx.BodyParser(r); if err != nil {
		return
	}

	if len(r.VersionID) == 0 {
		return errors.New("Provide version id")
	}

	if len(r.Prompt) == 0 {
		return errors.New("Provide prompt")
	} else if len(r.Prompt) > 1000 {
		return errors.New("Maximum prompt length is 1000")
	}

	if r.AmountImages != 1 && r.AmountImages != 4 {
		return errors.New("You can generate only 1 or 4 images")
	}

	return nil
}
