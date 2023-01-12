package payments

import (
	"ecomdream/src/domain/models"
	"ecomdream/src/pkg/external/informer"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v73"
)

func checkoutSessionCompleted(ctx *fiber.Ctx, event stripe.Event) error {
	var stripeSession stripe.CheckoutSession

	if err := json.Unmarshal(event.Data.Raw, &stripeSession); err != nil {
		logrus.Error(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	payment, err := models.GetPayment(stripeSession.ClientReferenceID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if payment == nil {
		logrus.Warningf("Unkow payment with id %s", stripeSession.ClientReferenceID)
		return ctx.SendStatus(fiber.StatusOK)
	}

	payment.Email = &stripeSession.CustomerDetails.Email
	payment.PaymentIntentID = &stripeSession.PaymentIntent.ID
	payment.SessionID = &stripeSession.ID
	payment.AmountPaid = &stripeSession.AmountTotal

	if err := payment.MarkAsPaid(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	go func() {
		logrus.Infof("Recieved payment %s", payment.ID)
		informer.SendTelegramMessage(fmt.Sprintf("Plan: %d, +%+v$", payment.PlanID, stripeSession.AmountTotal/100), "payments")
	}()

	return ctx.SendStatus(fiber.StatusOK)
}
