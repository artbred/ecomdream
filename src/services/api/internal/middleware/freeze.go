package middleware

import (
	"context"
	"ecomdream/src/pkg/storages/redisdb"
	"github.com/gofiber/fiber/v2"
	"time"
)

//TODO ctx.Params("id") does not work because it is init before endpoint

func FreezeEndpointForID(ctx *fiber.Ctx) error {
	id := string(ctx.Request().URI().LastPathSegment()); if len(id) == 0 {
		return ctx.Next()
	}

	key := redisdb.BuildFreezeEndpointKey(string(ctx.Request().URI().Path()), id)
	rdb := redisdb.Connection()

	if key.IsFrozen() {
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
