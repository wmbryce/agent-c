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

// TransformRequest converts the standard request format to provider-specific format using RequestSchema
func TransformRequest(modelKey string, messages []types.ChatMessage, options map[string]interface{}, config *types.ProviderConfig) map[string]interface{} {
	// Default schema (OpenAI-compatible)
	schema := &types.RequestSchema{
		ModelField:    "model",
		MessagesField: "messages",
	}

	if config != nil && config.RequestSchema != nil {
		schema = config.RequestSchema
	}

	request := make(map[string]interface{})

	// Add model field if specified
	if schema.ModelField != "" {
		request[schema.ModelField] = modelKey
	}

	// Transform and add messages
	transformedMessages := transformMessages(messages, schema.MessageTransform)
	if schema.MessagesField != "" {
		request[schema.MessagesField] = transformedMessages
	}

	// Transform and add options
	transformedOptions := transformOptions(options, schema)
	if schema.OptionsWrapper != "" {
		// Wrap options in a nested object
		if len(transformedOptions) > 0 {
			request[schema.OptionsWrapper] = transformedOptions
		}
	} else {
		// Add options directly to request
		for k, v := range transformedOptions {
			request[k] = v
		}
	}

	return request
}

// transformMessages applies MessageTransform to each message
func transformMessages(messages []types.ChatMessage, transform *types.MessageTransform) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(messages))

	// Default transform settings
	roleField := "role"
	contentPath := "content"
	var roleMap map[string]string

	if transform != nil {
		if transform.RoleField != "" {
			roleField = transform.RoleField
		}
		if transform.ContentPath != "" {
			contentPath = transform.ContentPath
		}
		roleMap = transform.RoleMap
	}

	for _, msg := range messages {
		role := msg.Role
		if roleMap != nil {
			if mappedRole, ok := roleMap[role]; ok {
				role = mappedRole
			}
		}

		message := make(map[string]interface{})
		message[roleField] = role

		// Handle content path (e.g., "content" or "parts[].text")
		setContentByPath(message, contentPath, msg.Content)

		result = append(result, message)
	}

	return result
}

// setContentByPath sets content at the specified path
// Supports simple paths ("content") and array paths ("parts[].text")
func setContentByPath(message map[string]interface{}, path string, content string) {
	// Check for array notation: "parts[].text"
	if idx := indexOf(path, "[]."); idx != -1 {
		arrayField := path[:idx]
		nestedField := path[idx+3:]

		message[arrayField] = []map[string]interface{}{
			{nestedField: content},
		}
		return
	}

	// Simple field
	message[path] = content
}

// indexOf returns the index of substr in s, or -1 if not found
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// transformOptions applies option renaming and filtering
func transformOptions(options map[string]interface{}, schema *types.RequestSchema) map[string]interface{} {
	if options == nil {
		return nil
	}

	// Build set of options to omit
	omitSet := make(map[string]bool)
	for _, opt := range schema.OptionsOmit {
		omitSet[opt] = true
	}

	result := make(map[string]interface{})
	for k, v := range options {
		// Skip omitted options
		if omitSet[k] {
			continue
		}

		// Rename if mapping exists
		newKey := k
		if schema.OptionsRename != nil {
			if renamed, ok := schema.OptionsRename[k]; ok {
				newKey = renamed
			}
		}

		result[newKey] = v
	}

	return result
}
