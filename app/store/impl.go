package store

import (
	"context"

	"github.com/wmbryce/agent-c/app/store/postgres"
	"github.com/wmbryce/agent-c/app/types"
)

type SqlStore interface {
	CreateModel(model *types.Model) (*types.Model, error)
	GetModels() ([]types.Model, error)
	GetModelCredentials(modelKey string) (*types.ModelCredentials, error)
	Close()
}

func NewSqlStore(ctx context.Context) SqlStore {
	return postgres.New(ctx)
}
