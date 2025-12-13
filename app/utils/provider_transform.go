package utils

import (
	"encoding/json"
	"net/http"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/wmbryce/agent-c/app/types"
)

// SetProviderHeaders sets the authentication and extra headers based on provider config
func SetProviderHeaders(req *http.Request, config *types.ProviderConfig, apiKey string) {
	if config == nil {
		// Default to Bearer token if no config
		req.Header.Set("Authorization", "Bearer "+apiKey)
		return
	}

	// Set authentication header
	switch config.AuthType {
	case "api_key":
		req.Header.Set(config.AuthHeader, apiKey)
	case "bearer":
		fallthrough
	default:
		if config.AuthHeader != "" {
			req.Header.Set(config.AuthHeader, "Bearer "+apiKey)
		} else {
			req.Header.Set("Authorization", "Bearer "+apiKey)
		}
	}

	// Set extra headers
	for key, value := range config.ExtraHeaders {
		req.Header.Set(key, value)
	}
}

// ApplyRequestDefaults merges default values into the request payload
// Only sets defaults for keys that are not already present
func ApplyRequestDefaults(request map[string]interface{}, defaults map[string]any) {
	for key, value := range defaults {
		if _, exists := request[key]; !exists {
			request[key] = value
		}
	}
}

// TransformResponse uses JSONPath mappings to transform a provider response into GeneralChatResponse
func TransformResponse(body []byte, mapping map[string]string) (*types.GeneralChatResponse, error) {
	// Parse the JSON response
	obj, err := oj.Parse(body)
	if err != nil {
		return nil, err
	}

	response := &types.GeneralChatResponse{}

	// Extract each field using JSONPath
	for field, path := range mapping {
		if path == "" || path == "null" {
			continue
		}

		expr, err := jp.ParseString(path)
		if err != nil {
			continue
		}

		results := expr.Get(obj)
		if len(results) == 0 {
			continue
		}

		value := results[0]

		switch field {
		case "id":
			if v, ok := value.(string); ok {
				response.ID = v
			}
		case "model":
			if v, ok := value.(string); ok {
				response.Model = v
			}
		case "content":
			if v, ok := value.(string); ok {
				response.Content = v
			}
		case "role":
			if v, ok := value.(string); ok {
				response.Role = v
			}
		case "finish_reason":
			if v, ok := value.(string); ok {
				response.FinishReason = v
			}
		case "prompt_tokens":
			response.PromptTokens = toInt(value)
		case "completion_tokens":
			response.CompletionTokens = toInt(value)
		case "total_tokens":
			response.TotalTokens = toInt(value)
		}
	}

	// Calculate total_tokens if not provided
	if response.TotalTokens == 0 && (response.PromptTokens > 0 || response.CompletionTokens > 0) {
		response.TotalTokens = response.PromptTokens + response.CompletionTokens
	}

	return response, nil
}

// toInt converts various numeric types to int
func toInt(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case json.Number:
		if i, err := v.Int64(); err == nil {
			return int(i)
		}
	}
	return 0
}
