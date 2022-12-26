package middleware

import (
	"context"
	"ecomdream/src/pkg/storages/redisdb"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"time"
)

func buildFreezeRedisKey(endpoint, id string) string {
	return fmt.Sprintf("freeze_endpoint:%s,%s", endpoint, id)
}

func FreezeEndpointForID(ctx *fiber.Ctx) error {
	id := ctx.Query("id")
	if len(id) == 0 {
		return ctx.Next()
	}

	key := buildFreezeRedisKey(string(ctx.Request().URI().Path()), id)
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
			"message": "Please wait for previous request to finish",
		})
	} else {
		rdb.SetNX(context.Background(), key, true, 10*time.Second)
	}

	err = ctx.Next()
	rdb.Del(context.Background(), key)
	return err
}
