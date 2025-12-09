package types

import "time"

type ConsumeModelRequest struct {
	ModelKey string                 `json:"model_key" validate:"required"`
	Messages []ChatMessage          `json:"messages" validate:"required,min=1,dive"`
	Options  map[string]interface{} `json:"options,omitempty"`
	MaxCost  float64                `json:"max_cost" validate:"required,gt=0"`
}

type ModelCredentials struct {
	ModelKey        string `json:"model_key"`
	RequestURL      string `json:"request_url"`
	ApiKey          string `json:"api_key"`
	TokensAvailable int    `json:"tokens_available"`
	ProviderName    string `json:"provider_name"`
}

type Model struct {
	ID               string     `json:"id"`
	ModelKey         string     `json:"model_key" validate:"required"`
	Name             string     `json:"name" validate:"required"`
	Description      string     `json:"description" validate:"required"`
	ProviderID       string     `json:"provider_id" validate:"required,uuid"`
	OptionsSchemaID  string     `json:"options_schema_id" validate:"required,uuid"`
	ResponseSchemaID string     `json:"response_schema_id" validate:"required,uuid"`
	RequestURL       string     `json:"request_url" validate:"required,url"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at" db:"updated_at"`
}
