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
	RequestSchema   *RequestSchema    `json:"request_schema"`
}

// RequestSchema defines how to transform standard requests to provider-specific format
type RequestSchema struct {
	// ModelField is the field name for the model key (e.g., "model"). Empty/null to omit.
	ModelField string `json:"model_field,omitempty"`

	// MessagesField is the field name for the messages array (e.g., "messages", "contents")
	MessagesField string `json:"messages_field"`

	// MessageTransform defines how to transform each message
	MessageTransform *MessageTransform `json:"message_transform,omitempty"`

	// OptionsWrapper wraps all options in a nested object (e.g., "generationConfig")
	OptionsWrapper string `json:"options_wrapper,omitempty"`

	// OptionsRename maps standard option names to provider-specific names
	OptionsRename map[string]string `json:"options_rename,omitempty"`

	// OptionsOmit lists option names to exclude from the request
	OptionsOmit []string `json:"options_omit,omitempty"`
}

// MessageTransform defines how to transform message structure
type MessageTransform struct {
	// RoleField is the field name for the role (default: "role")
	RoleField string `json:"role_field,omitempty"`

	// ContentPath is the path where content goes (e.g., "content" or "parts[].text")
	ContentPath string `json:"content_path,omitempty"`

	// RoleMap remaps role values (e.g., {"assistant": "model"})
	RoleMap map[string]string `json:"role_map,omitempty"`
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
