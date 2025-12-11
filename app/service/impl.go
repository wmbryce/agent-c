package service

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/wmbryce/agent-c/app/store"
)

// HTTPClient interface for making HTTP requests (allows mocking in tests)
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Service struct {
	logger     *zerolog.Logger
	store      store.SqlStore
	fiber      *fiber.App
	httpClient HTTPClient
}

func New(logger *zerolog.Logger, sqlStore store.SqlStore, fiber *fiber.App, client HTTPClient) *Service {
	if client == nil {
		client = &http.Client{}
	}
	return &Service{
		logger:     logger,
		store:      sqlStore,
		fiber:      fiber,
		httpClient: client,
	}
}
