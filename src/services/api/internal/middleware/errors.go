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

	response := ctx.Response()
	if response != nil {
		statusCode := response.StatusCode()
		if statusCode < 200 || statusCode >= 300 {
			endpoint := string(ctx.Request().URI().Path())
			go func() {
				if statusCode == 500 {
					informer.SendTelegramMessage(fmt.Sprintf("%s\n%s", time.Now().Format(time.RFC3339), endpoint), informer.InternalLevel)
				}
				monitoring.EndpointErrors.WithLabelValues(strconv.Itoa(statusCode), endpoint).Inc()
			}()
		}
	}

	return err
}
