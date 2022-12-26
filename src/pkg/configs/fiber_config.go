package configs

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
	"time"
)

func LoadFiberServerConfig() fiber.Config {
	return fiber.Config{
		ReadTimeout: time.Second * time.Duration(60),
		BodyLimit: 6 * 1024 * 1024,
	}
}
