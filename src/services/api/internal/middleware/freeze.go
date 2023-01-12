package middleware

import (
	"context"
	"ecomdream/src/pkg/storages/redisdb"
	"github.com/gofiber/fiber/v2"
	"time"
)

func FreezeEndpointForID(ctx *fiber.Ctx) error {
	id := string(ctx.Request().URI().LastPathSegment()); if len(id) == 0 { //TODO ctx.Params("id")
		return ctx.Next()
	}

	key := redisdb.BuildFreezeEndpointKey(string(ctx.Request().URI().Path()), id)
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
