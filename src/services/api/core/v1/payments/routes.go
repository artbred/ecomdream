package payments

import "github.com/gofiber/fiber/v2"

func Init(router fiber.Router) {
	router.Post("/create", CreatePaymentLinkHandler)
	router.Post("/webhook", WebhookListenerHandler)
}
