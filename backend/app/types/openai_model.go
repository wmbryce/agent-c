package types

// ChatCompletionRequest struct to describe chat completion request object.
type OpenAIChatCompletionRequest struct {
	Model       string          `json:"model" validate:"required"`
	Messages    []ChatMessage   `json:"messages" validate:"required,min=1,dive"`
	Temperature *float32        `json:"temperature,omitempty" validate:"omitempty,min=0,max=2"`
	MaxTokens   *int            `json:"max_tokens,omitempty" validate:"omitempty,min=1"`
	Stream      bool            `json:"stream,omitempty"`
}

// ChatMessage struct to describe chat message object.
type ChatMessage struct {
	Role    string `json:"role" validate:"required,oneof=system user assistant"`
	Content string `json:"content" validate:"required"`
}

// ChatCompletionResponse struct to describe chat completion response object.
type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice struct to describe choice object.
type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// Usage struct to describe usage object.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

