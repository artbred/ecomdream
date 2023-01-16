package prompts

import (
	"ecomdream/src/services/api/internal/middleware"
	"github.com/gofiber/fiber/v2"
)


func Init(router fiber.Router) {
	h := createHandler()

	router.Use(middleware.FreezeEndpointForID).Post("/create/:id", h.CreatePromptForIDHandler)
	router.Get("/list/:id", h.ListPromptsForIDHandler)
}
