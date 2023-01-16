package server

import (
	"ecomdream/src/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
)

func StartServerWithGracefulShutdown(a *fiber.App) {
	idleConnClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := a.Shutdown(); err != nil {
			logrus.Infof("Oops... Server is not shutting down! Reason: %v", err)
		}

		close(idleConnClosed)
	}()

	fiberConnURL, _ := config.ConnectionURLBuilder("fiber")

	if err := a.Listen(fiberConnURL); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnClosed
}

func StartServer(a *fiber.App) {
	fiberConnURL, _ := config.ConnectionURLBuilder("fiber")

	if err := a.Listen(fiberConnURL); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}
}
