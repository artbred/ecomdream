package versions

import (
	"ecomdream/src/services/api/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func Init(router fiber.Router) {
	h := createHandler()

	router.Use(middleware.FreezeEndpointForID).Post("/train/:id", h.TrainVersionHandler)
	router.Get("/info/:id", h.VersionInfoHandler)
}
