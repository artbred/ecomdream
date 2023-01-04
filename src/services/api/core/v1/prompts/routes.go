package prompts

import (
	"ecomdream/src/services/api/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func Init(router fiber.Router) {
	router.Use(middleware.FreezeEndpointForID).Post("/create/:id", CreatePromptForIDHandler)
	router.Get("/list/:id", ListPromptsForIDHandler)
}
