package types

import "time"

type ConsumeModelRequest struct {
	ModelKey string                 `json:"model_key" validate:"required"`
	Messages []ChatMessage          `json:"messages" validate:"required,min=1,dive"`
	Options  map[string]interface{} `json:"options,omitempty"`
	MaxCost  float64                `json:"max_cost" validate:"required,gt=0"`
}

type ModelCredentials struct {
	ModelKey        string          `json:"model_key"`
	RequestURL      string          `json:"request_url"`
	ApiKey          string          `json:"api_key"`
	TokensAvailable int             `json:"tokens_available"`
	ProviderName    string          `json:"provider_name"`
	ProviderConfig  *ProviderConfig `json:"provider_config"`
}

type ProviderConfig struct {
	AuthType        string            `json:"auth_type"`
	AuthHeader      string            `json:"auth_header"`
	ExtraHeaders    map[string]string `json:"extra_headers"`
	RequestDefaults map[string]any    `json:"request_defaults"`
	ResponseMapping map[string]string `json:"response_mapping"`
}

type GeneralChatResponse struct {
	ID               string `json:"id"`
	Model            string `json:"model"`
	Content          string `json:"content"`
	Role             string `json:"role"`
	FinishReason     string `json:"finish_reason"`
	PromptTokens     int    `json:"prompt_tokens"`
	CompletionTokens int    `json:"completion_tokens"`
	TotalTokens      int    `json:"total_tokens"`
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
