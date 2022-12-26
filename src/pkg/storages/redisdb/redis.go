package redisdb

import (
	"context"
	"ecomdream/src/pkg/configs"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"sync"
)

var connection *redis.Client
var once sync.Once

func Connection() *redis.Client {
	if connection == nil {
		once.Do(Init)
	}

	return connection
}

func Init() {
	redisAddr, _ := configs.ConnectionURLBuilder("redis")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	status := rdb.Ping(context.Background())
	if status == nil {
		logrus.Errorf("Can't connect to redis")
		return
	}

	if status.String() != "ping: PONG" {
		logrus.Errorf("Can't ping redis, %s", status.String())
		return
	}

	connection = rdb
}

func init() {
	Init()
}
