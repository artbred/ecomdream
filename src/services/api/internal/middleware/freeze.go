package middleware

import (
	"context"
	"ecomdream/src/pkg/storages/redisdb"
	"github.com/gofiber/fiber/v2"
	"time"
)

func FreezeEndpointForID(ctx *fiber.Ctx) error {
	id := ctx.Query("id"); if len(id) == 0 {
		return ctx.Next()
	}

	key := redisdb.BuildBlockEndpointKey(string(ctx.Request().URI().Path()), id)
	rdb := redisdb.Connection()

	if key.IsBlocked() {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":    fiber.StatusForbidden,
			"message": "You have running task, please wait for it to complete, maximum wait time is 5 minutes",
		})
	} else {
		rdb.SetNX(context.Background(), string(key), true, 5*time.Minute)
	}

	err := ctx.Next()
	rdb.Del(context.Background(), string(key))
	return err
}
