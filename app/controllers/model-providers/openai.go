package controllers

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	openai "github.com/sashabaranov/go-openai"
	"github.com/wmbryce/agent-c/app/types"
	"github.com/wmbryce/agent-c/app/utils"
)

// ChatCompletion func sends a chat completion request to OpenAI API.
// @Description Send a chat completion request to OpenAI API.
// @Summary send a chat completion request to OpenAI
// @Tags OpenAI
// @Accept json
// @Produce json
// @Param request body types.ChatCompletionRequest true "Chat completion request"
// @Success 200 {object} types.ChatCompletionResponse
// @Security ApiKeyAuth
// @Router /v1/openai/chat [post]
func ChatCompletion(c *fiber.Ctx) error {
	// Create new ChatCompletionRequest struct
	request := &types.ChatCompletionRequest{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(request); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new validator for a ChatCompletionRequest model.
	validate := utils.NewValidator()

	// Validate request fields.
	if err := validate.Struct(request); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Get OpenAI API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "OPENAI_API_KEY environment variable is not set",
		})
	}

	// Create OpenAI client
	client := openai.NewClient(apiKey)

	// Build messages for OpenAI API
	messages := make([]openai.ChatCompletionMessage, len(request.Messages))
	for i, msg := range request.Messages {
		messages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Build OpenAI request
	openaiReq := openai.ChatCompletionRequest{
		Model:    request.ModelKey,
		Messages: messages,
	}

	// Add optional parameters
	if request.Options["temperature"] != nil {
		openaiReq.Temperature = request.Options["temperature"].(float32)
	}
	if request.Options["max_tokens"] != nil {
		openaiReq.MaxCompletionTokens = request.Options["max_tokens"].(int)
	}

	// Send request to OpenAI API
	resp, err := client.CreateChatCompletion(context.Background(), openaiReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Build response
	choices := make([]types.Choice, len(resp.Choices))
	for i := range resp.Choices {
		choices[i] = types.Choice{
			Index: resp.Choices[i].Index,
			Message: types.ChatMessage{
				Role:    resp.Choices[i].Message.Role,
				Content: resp.Choices[i].Message.Content,
			},
			FinishReason: string(resp.Choices[i].FinishReason),
		}
	}

	response := types.ChatCompletionResponse{
		ID:      resp.ID,
		Object:  resp.Object,
		Created: resp.Created,
		Model:   resp.Model,
		Choices: choices,
		Usage: types.Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"response": response,
	})
}

