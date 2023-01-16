package middleware

import (
	"ecomdream/src/pkg/external/informer"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

func ErrorMiddleware(ctx *fiber.Ctx) error {
	err := ctx.Next()

	statusCode := ctx.Response().StatusCode()

	if statusCode == fiber.StatusInternalServerError {
		go informer.SendTelegramMessage(fmt.Sprintf("%s\n%s", time.Now().Format(time.RFC3339), string(ctx.Request().URI().Path())), informer.InternalLevel)
	}



	return err
}
