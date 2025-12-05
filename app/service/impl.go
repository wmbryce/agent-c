package service

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/wmbryce/agent-c/app/store"
	"github.com/wmbryce/agent-c/app/store/cache"
)

type Service struct {
	logger *log.Logger
	store *store.Store
	cache *cache.Cache
	fiber *fiber.App
}

func New(logger *log.Logger, store *store.Store, cache *cache.Cache, fiber *fiber.App) *Service {
	return &Service{logger: logger, store: store}
}