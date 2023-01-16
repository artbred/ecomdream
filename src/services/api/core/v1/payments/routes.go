package payments

import "github.com/gofiber/fiber/v2"

func Init(router fiber.Router) {
	h := createHandler()

	router.Post("/create", h.CreatePaymentLinkHandler)
	router.Post("/webhook", h.WebhookListenerHandler)
}
