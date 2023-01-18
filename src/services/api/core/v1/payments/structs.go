package payments

import (
	"ecomdream/src/domain/models"
	"github.com/gofiber/fiber/v2"
)

type CreatePaymentLinkRequest struct {
	PlanID int `json:"plan_id"`
	VersionID *string `json:"version_id,omitempty"`
	PromocodeID *string `json:"promocode_id,omitempty"`
}

type AvailablePlansResponse struct {
	Code int `json:"code"`
	Plans []models.Plan `json:"plans"`
}

func (r *CreatePaymentLinkRequest) Validate(ctx *fiber.Ctx) (err error) {
	err = ctx.BodyParser(r); if err != nil {
		return
	}

	if r.VersionID != nil {
		if len(*r.VersionID) == 0 {
			r.VersionID = nil
		}
	}

	if r.PromocodeID != nil {
		if len(*r.PromocodeID) == 0 {
			r.VersionID = nil
		}
	}

	return nil
}

type CreatePaymentLinkResponse struct {
	Code int `json:"code"`
	URL string `json:"url"`
}
