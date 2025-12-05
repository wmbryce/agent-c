package utils

import (
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func StartServerWithGracefulShutdown(a *fiber.App) {
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := a.Shutdown(); err != nil {
			log.Error().Err(err).Msg("server is not shutting down")
		}

		close(idleConnsClosed)
	}()

	fiberConnURL, _ := ConnectionURLBuilder("fiber")

	if err := a.Listen(fiberConnURL); err != nil {
		log.Error().Err(err).Msg("server is not running")
	}

	<-idleConnsClosed
}

func StartServer(a *fiber.App) {
	fiberConnURL, _ := ConnectionURLBuilder("fiber")

	if err := a.Listen(fiberConnURL); err != nil {
		log.Error().Err(err).Msg("server is not running")
	}
}
