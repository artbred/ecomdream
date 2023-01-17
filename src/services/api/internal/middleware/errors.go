package middleware

import (
	"ecomdream/src/pkg/external/informer"
	"ecomdream/src/services/api/internal/monitoring"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

func ErrorMiddleware(ctx *fiber.Ctx) error {
	err := ctx.Next()

	go func(statusCode int) {
		if statusCode < 200 || statusCode >= 300 {
			endpoint := string(ctx.Request().URI().Path())
			if statusCode == 500 {
				go informer.SendTelegramMessage(fmt.Sprintf("%s\n%s", time.Now().Format(time.RFC3339), endpoint), informer.InternalLevel)
			}
			monitoring.EndpointErrors.WithLabelValues(strconv.Itoa(statusCode), endpoint).Inc()
		}
	}(ctx.Response().StatusCode())

	return err
}
