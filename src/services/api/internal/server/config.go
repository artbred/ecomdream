package server

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

func LoadFiberServerConfig() fiber.Config {
	return fiber.Config{
		ReadTimeout: time.Second * time.Duration(60),
		BodyLimit: 6 * 1024 * 1024,
	}
}
