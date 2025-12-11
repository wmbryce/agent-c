package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/wmbryce/agent-c/app/service"
	"github.com/wmbryce/agent-c/app/types"
)

func TestConsumeModel(t *testing.T) {
	logger := zerolog.Nop()

	tests := []struct {
		name           string
		requestBody    interface{}
		mockCreds      *types.ModelCredentials
		mockCredsErr   error
		mockHTTPResp   *http.Response
		mockHTTPErr    error
		expectedStatus int
		expectedError  bool
		expectedMsg    string
	}{
		{
			name: "successful request",
			requestBody: types.ConsumeModelRequest{
				ModelKey: "gpt-4",
				Messages: []types.ChatMessage{
					{Role: "user", Content: "Hello"},
				},
				MaxCost: 100,
			},
			mockCreds: &types.ModelCredentials{
				ModelKey:        "gpt-4",
				RequestURL:      "https://api.openai.com/v1/chat/completions",
				ApiKey:          "sk-test-key",
				TokensAvailable: 1000,
				ProviderName:    "openai",
			},
			mockHTTPResp: &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(strings.NewReader(`{
					"id": "chatcmpl-123",
					"object": "chat.completion",
					"created": 1234567890,
					"model": "gpt-4",
					"choices": [{"index": 0, "message": {"role": "assistant", "content": "Hi there!"}, "finish_reason": "stop"}],
					"usage": {"prompt_tokens": 10, "completion_tokens": 20, "total_tokens": 30}
				}`)),
			},
			expectedStatus: 200,
			expectedError:  false,
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			expectedStatus: 400,
			expectedError:  true,
		},
		{
			name: "missing required fields",
			requestBody: map[string]interface{}{
				"model_key": "gpt-4",
				// missing messages and max_cost
			},
			expectedStatus: 400,
			expectedError:  true,
		},
		{
			name: "model not found",
			requestBody: types.ConsumeModelRequest{
				ModelKey: "non-existent-model",
				Messages: []types.ChatMessage{
					{Role: "user", Content: "Hello"},
				},
				MaxCost: 100,
			},
			mockCredsErr:   errors.New("model not found"),
			expectedStatus: 404,
			expectedError:  true,
			expectedMsg:    "model not found or no API key available",
		},
		{
			name: "insufficient tokens",
			requestBody: types.ConsumeModelRequest{
				ModelKey: "gpt-4",
				Messages: []types.ChatMessage{
					{Role: "user", Content: "Hello"},
				},
				MaxCost: 1000,
			},
			mockCreds: &types.ModelCredentials{
				ModelKey:        "gpt-4",
				RequestURL:      "https://api.openai.com/v1/chat/completions",
				ApiKey:          "sk-test-key",
				TokensAvailable: 100, // less than MaxCost
				ProviderName:    "openai",
			},
			expectedStatus: 402,
			expectedError:  true,
			expectedMsg:    "insufficient tokens available",
		},
		{
			name: "provider unreachable",
			requestBody: types.ConsumeModelRequest{
				ModelKey: "gpt-4",
				Messages: []types.ChatMessage{
					{Role: "user", Content: "Hello"},
				},
				MaxCost: 100,
			},
			mockCreds: &types.ModelCredentials{
				ModelKey:        "gpt-4",
				RequestURL:      "https://api.openai.com/v1/chat/completions",
				ApiKey:          "sk-test-key",
				TokensAvailable: 1000,
				ProviderName:    "openai",
			},
			mockHTTPErr:    errors.New("connection refused"),
			expectedStatus: 502,
			expectedError:  true,
			expectedMsg:    "failed to reach model provider",
		},
		{
			name: "invalid provider response",
			requestBody: types.ConsumeModelRequest{
				ModelKey: "gpt-4",
				Messages: []types.ChatMessage{
					{Role: "user", Content: "Hello"},
				},
				MaxCost: 100,
			},
			mockCreds: &types.ModelCredentials{
				ModelKey:        "gpt-4",
				RequestURL:      "https://api.openai.com/v1/chat/completions",
				ApiKey:          "sk-test-key",
				TokensAvailable: 1000,
				ProviderName:    "openai",
			},
			mockHTTPResp: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`invalid json`)),
			},
			expectedStatus: 500,
			expectedError:  true,
			expectedMsg:    "failed to parse provider response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock store
			store := &MockStore{
				Creds:    tt.mockCreds,
				CredsErr: tt.mockCredsErr,
			}

			// Create mock HTTP client
			httpClient := &MockHTTPClient{
				Response: tt.mockHTTPResp,
				Err:      tt.mockHTTPErr,
			}

			// Create Fiber app and service
			app := fiber.New()
			svc := service.New(&logger, store, app, httpClient)

			// Register route
			app.Post("/api/v1/ai/consume", svc.ConsumeModel)

			// Create request body
			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
			}

			// Create test request
			req := httptest.NewRequest("POST", "/api/v1/ai/consume", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("failed to execute request: %v", err)
			}

			// Check status code
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			// Parse response
			var result map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			// Check error field
			if errorVal, ok := result["error"].(bool); ok {
				if errorVal != tt.expectedError {
					t.Errorf("expected error=%v, got %v", tt.expectedError, errorVal)
				}
			}

			// Check message if expected
			if tt.expectedMsg != "" {
				if msg, ok := result["msg"].(string); ok {
					if msg != tt.expectedMsg {
						t.Errorf("expected msg=%q, got %q", tt.expectedMsg, msg)
					}
				}
			}
		})
	}
}
