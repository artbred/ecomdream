package versions

import (
	"ecomdream/src/services/api/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func Init(router fiber.Router) {
	router.Use(middleware.FreezeEndpointForID).Post("/train/:id", TrainVersionHandler)
	router.Get("/info/:id", VersionInfoHandler)
}
