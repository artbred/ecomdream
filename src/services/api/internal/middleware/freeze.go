package middleware

import (
	"context"
	"ecomdream/src/pkg/storages/redisdb"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"time"
)

func FreezeEndpointForID(ctx *fiber.Ctx) error {
	id := ctx.Query("id"); if len(id) == 0 {
		return ctx.Next()
	}

	key := redisdb.BuildFreezeEndpointKey(string(ctx.Request().URI().Path()), id)
	rdb := redisdb.Connection()

	isBlocked, err := rdb.Get(context.Background(), key).Bool()
	if err != nil {
		if err == redis.Nil {
			isBlocked = false
		} else {
			logrus.Error(err)
			return ctx.Next()
		}
	}

	if isBlocked {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":    fiber.StatusForbidden,
			"message": "You have running task, please wait for it to complete, maximum wait time is 5 minutes",
		})
	} else {
		rdb.SetNX(context.Background(), key, true, 5*time.Minute)
	}

	err = ctx.Next()
	rdb.Del(context.Background(), key)
	return err
}
