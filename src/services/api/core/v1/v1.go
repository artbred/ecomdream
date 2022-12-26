package v1

import (
	"ecomdream/src/services/api/core/v1/payments"
	"ecomdream/src/services/api/core/v1/prompts"
	"ecomdream/src/services/api/core/v1/versions"
	"github.com/gofiber/fiber/v2"
)

func Init(router fiber.Router) {
	router = router.Group("/v1")

	payments.Init(router.Group("/payments"))
	versions.Init(router.Group("/versions"))
	prompts.Init(router.Group("/prompts"))
}
