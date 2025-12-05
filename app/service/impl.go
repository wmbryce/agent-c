package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/wmbryce/agent-c/app/store"
)

type Service struct {
	logger *zerolog.Logger
	store  store.SqlStore
	fiber  *fiber.App
}

func New(logger *zerolog.Logger, sqlStore store.SqlStore, fiber *fiber.App) *Service {
	return &Service{logger: logger, store: sqlStore, fiber: fiber}
}
