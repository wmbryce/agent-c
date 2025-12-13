package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wmbryce/agent-c/app/types"
)

func (s *Store) CreateModel(model *types.Model) (*types.Model, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if model.ID == "" {
		model.ID = uuid.New().String()
	}

	query := `
		INSERT INTO agc.models (id, model_key, name, description, provider_id, options_schema_id, response_schema_id, request_url, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, model_key, request_url, created_at, updated_at
	`

	var createdModel types.Model
	err := s.db.QueryRow(ctx, query,
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

	createdModel.Name = model.Name
	createdModel.Description = model.Description
	createdModel.ProviderID = model.ProviderID
	createdModel.OptionsSchemaID = model.OptionsSchemaID
	createdModel.ResponseSchemaID = model.ResponseSchemaID

	return &createdModel, nil
}

func (s *Store) GetModels() ([]types.Model, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, model_key, name, description, provider_id, options_schema_id, response_schema_id, request_url, created_at, updated_at
		FROM agc.models
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(ctx, query)
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

func (s *Store) GetModelCredentials(modelKey string) (*types.ModelCredentials, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT m.model_key, m.request_url, ak.api_key, ak.tokens_available, p.name,
		       p.auth_type, p.auth_header, p.extra_headers, p.request_defaults, p.response_mapping
		FROM agc.models m
		JOIN agc.providers p ON m.provider_id = p.id
		JOIN agc.api_keys ak ON ak.provider_id = p.id
		WHERE m.model_key = $1
	`

	var creds types.ModelCredentials
	var authType, authHeader *string
	var extraHeaders, requestDefaults, responseMapping []byte

	err := s.db.QueryRow(ctx, query, modelKey).Scan(
		&creds.ModelKey,
		&creds.RequestURL,
		&creds.ApiKey,
		&creds.TokensAvailable,
		&creds.ProviderName,
		&authType,
		&authHeader,
		&extraHeaders,
		&requestDefaults,
		&responseMapping,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get model credentials: %w", err)
	}

	// Parse provider config
	config := &types.ProviderConfig{
		ExtraHeaders:    make(map[string]string),
		RequestDefaults: make(map[string]any),
		ResponseMapping: make(map[string]string),
	}

	if authType != nil {
		config.AuthType = *authType
	}
	if authHeader != nil {
		config.AuthHeader = *authHeader
	}
	if len(extraHeaders) > 0 {
		json.Unmarshal(extraHeaders, &config.ExtraHeaders)
	}
	if len(requestDefaults) > 0 {
		json.Unmarshal(requestDefaults, &config.RequestDefaults)
	}
	if len(responseMapping) > 0 {
		json.Unmarshal(responseMapping, &config.ResponseMapping)
	}

	creds.ProviderConfig = config

	return &creds, nil
}
