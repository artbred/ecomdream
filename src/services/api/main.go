package main

import (
	"ecomdream/src/pkg/config"
	v1 "ecomdream/src/services/api/core/v1"
	"ecomdream/src/services/api/internal/middleware"
	"ecomdream/src/services/api/internal/server"
	"ecomdream/src/services/api/internal/swagger"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
	"strconv"
)

func SetupAPI(router fiber.Router) {
	v1.Init(router)
}

// @title API
// @version 1.0
// @BasePath /api
func main() {
	app := fiber.New(server.LoadFiberServerConfig())

	middleware.SetupMiddleware(app)
	swagger.ServeDocs(app)

	app.Get("/health", func(ctx *fiber.Ctx) error {
		return nil
	})

	SetupAPI(app.Group("/api"))

	if config.Debug {
		server.StartServer(app)
	} else {
		server.StartServerWithGracefulShutdown(app)
	}
}

func init() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
			return "", fmt.Sprintf("[%s]", fileName)
		},
		FullTimestamp: true,
	})
}
