package types

import "time"

type ChatCompletionRequest struct {
	ModelKey       string          `json:"model_key" validate:"required"`
	Messages       []ChatMessage   `json:"messages" validate:"required,min=1,dive"`
	Options        map[string]interface{} `json:"options,omitempty"`
}

type Model struct {
	ID       string          `json:"id"`
	ModelKey string          `json:"model_key"`
	Options  map[string]interface{} `json:"options"`
	ResponseSchema map[string]interface{} `json:"response_schema"`
	RequestURL string          `json:"request_url"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}