package payments

import (
	"ecomdream/src/domain/models"
	"ecomdream/src/pkg/configs"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v73"
	"github.com/stripe/stripe-go/v73/checkout/session"
	"github.com/stripe/stripe-go/v73/webhook"
	"github.com/twinj/uuid"
)

// CreatePaymentLinkHandler handler that creates payment link
// @Description Create payment link
// @Summary Create payment link
// @Tags payments
// @Accept json
// @Produce json
// @Param payment_data body CreatePaymentLinkRequest true "Payment data"
// @Success 201 {object} CreatePaymentLinkResponse
// @Router /v1/payments/create [post]
func CreatePaymentLinkHandler(ctx *fiber.Ctx) error {
	req := &CreatePaymentLinkRequest{}

	if err := req.Validate(ctx); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Bad request",
		})
	}

	plan, err := models.GetPlan(req.PlanID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Please try again later",
		})
	}

	if plan == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid plan id",
		})
	}

	if plan.IsDeprecated {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Plan is deprecated",
		})
	}

	payment := &models.Payment{
		ID:     uuid.NewV4().String(),
		PlanID: plan.ID,
	}

	successURL := fmt.Sprintf("https://kek.com/complete/%s", payment.ID)
	var discounts []*stripe.CheckoutSessionDiscountParams

	if !plan.IsInit {
		if req.VersionID == nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "You must also provide version id",
			})
		}

		version, err := models.GetVersion(*req.VersionID)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "Please try again later",
			})
		}

		if version == nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "Invalid version id",
			})
		}

		if version.PushedAt == nil {
			version = nil
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Please wait for your version to be ready",
			})
		}

		payment.VersionID = &version.ID
		successURL = fmt.Sprintf("https://kek.com/versions/%s", version.ID)
	}

	if req.PromocodeID != nil {
		discounts = append(discounts, &stripe.CheckoutSessionDiscountParams{
			PromotionCode: req.PromocodeID,
		})

		payment.PromocodeID = req.PromocodeID
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),

		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(stripe.CurrencyUSD)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:        stripe.String(plan.PlanName),
						Description: stripe.String(plan.PlanDescription),
					},
					UnitAmount: stripe.Int64(plan.Price),
				},
				Quantity: stripe.Int64(1),
			},
		},

		Mode:              stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:        stripe.String(successURL),
		CancelURL:         stripe.String("https://kek.com"),
		ClientReferenceID: stripe.String(payment.ID),

		ConsentCollection: &stripe.CheckoutSessionConsentCollectionParams{
			//Promotions: stripe.String("auto"),
			TermsOfService: stripe.String("required"),
		},

		Discounts: discounts,

		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Description: stripe.String(payment.ID),
			//Remove to capture automatically
			CaptureMethod: stripe.String(string(stripe.PaymentIntentCaptureMethodManual)),
		},

		//AutomaticTax: stripe.Bool(true),
	}

	s, err := session.New(params); if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Please try again later",
		})
	}

	if err := payment.Create(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Please try again later",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(CreatePaymentLinkResponse{
		Code: fiber.StatusCreated,
		URL:  s.URL,
	})
}

func WebhookListenerHandler(ctx *fiber.Ctx) error {
	event, err := webhook.ConstructEvent(ctx.Request().Body(), string(ctx.Request().Header.Peek("Stripe-Signature")), configs.StripeWebhookSecret)
	if err != nil {
		logrus.Error(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	switch event.Type {
	case "checkout.session.completed":
		return checkoutSessionCompleted(ctx, event)
	//case "checkout.session.expired":
	//	return checkoutSessionExpired(ctx, event)
	default:
		return ctx.SendStatus(fiber.StatusOK)
	}
}

func init() {
	stripe.Key = configs.StripeSecretKey
}
