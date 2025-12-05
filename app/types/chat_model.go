package types

import "time"

type ChatCompletionRequest struct {
	ModelKey       string          `json:"model_key" validate:"required"`
	Messages       []ChatMessage   `json:"messages" validate:"required,min=1,dive"`
	Options        map[string]interface{} `json:"options,omitempty"`
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