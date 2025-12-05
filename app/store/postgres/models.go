package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wmbryce/agent-c/app/types"
)

// CreateModel inserts a new model into the database
func CreateModel(pool *pgxpool.Pool, model *types.Model) (*types.Model, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Generate UUID if not provided
	if model.ID == "" {
		model.ID = uuid.New().String()
	}

	query := `
		INSERT INTO agc.models (id, model_key, name, description, provider_id, options_schema_id, response_schema_id, request_url, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, model_key, request_url, created_at, updated_at
	`

	var createdModel types.Model
	err := pool.QueryRow(ctx, query,
		model.ID,
		model.ModelKey,
		model.Name,
		model.Description,
		model.ProviderID,
		model.OptionsSchemaID,
		model.ResponseSchemaID,
		model.RequestURL,
		time.Now(),
	).Scan(
		&createdModel.ID,
		&createdModel.ModelKey,
		&createdModel.RequestURL,
		&createdModel.CreatedAt,
		&createdModel.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create model: %w", err)
	}

	// Copy over the fields we inserted but didn't return
	createdModel.Name = model.Name
	createdModel.Description = model.Description
	createdModel.ProviderID = model.ProviderID
	createdModel.OptionsSchemaID = model.OptionsSchemaID
	createdModel.ResponseSchemaID = model.ResponseSchemaID

	return &createdModel, nil
}

// GetModels retrieves all models from the database
func GetModels(pool *pgxpool.Pool) ([]types.Model, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, model_key, name, description, provider_id, options_schema_id, response_schema_id, request_url, created_at, updated_at
		FROM agc.models
		ORDER BY created_at DESC
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query models: %w", err)
	}
	defer rows.Close()

	var models []types.Model
	for rows.Next() {
		var m types.Model
		err := rows.Scan(
			&m.ID,
			&m.ModelKey,
			&m.Name,
			&m.Description,
			&m.ProviderID,
			&m.OptionsSchemaID,
			&m.ResponseSchemaID,
			&m.RequestURL,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan model: %w", err)
		}
		models = append(models, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating models: %w", err)
	}

	return models, nil
}

