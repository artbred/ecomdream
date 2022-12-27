package middleware

import (
	"ecomdream/src/pkg/external/informer"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

func InformOnInternalError(ctx *fiber.Ctx) error {
	err := ctx.Next()

	if ctx.Response().StatusCode() == fiber.StatusInternalServerError {
		go informer.SendTelegramMessage(fmt.Sprintf("%s\n%s", time.Now().Format(time.RFC3339), string(ctx.Request().URI().Path())), "internal")
	}

	return err
}
