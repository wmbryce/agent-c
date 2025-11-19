package controllers

import (
	"context"
	"os"

	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/create-go-app/fiber-go-template/pkg/utils"
	"github.com/gofiber/fiber/v2"
	openai "github.com/sashabaranov/go-openai"
)

// ChatCompletion func sends a chat completion request to OpenAI API.
// @Description Send a chat completion request to OpenAI API.
// @Summary send a chat completion request to OpenAI
// @Tags OpenAI
// @Accept json
// @Produce json
// @Param request body models.ChatCompletionRequest true "Chat completion request"
// @Success 200 {object} models.ChatCompletionResponse
// @Security ApiKeyAuth
// @Router /v1/openai/chat [post]
func ChatCompletion(c *fiber.Ctx) error {
	// Create new ChatCompletionRequest struct
	request := &models.ChatCompletionRequest{}

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
		Model:    request.Model,
		Messages: messages,
	}

	// Add optional parameters
	if request.Temperature != nil {
		openaiReq.Temperature = *request.Temperature
	}
	if request.MaxTokens != nil {
		openaiReq.MaxTokens = *request.MaxTokens
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
	choices := make([]models.Choice, len(resp.Choices))
	for i, choice := range resp.Choices {
		choices[i] = models.Choice{
			Index: choice.Index,
			Message: models.ChatMessage{
				Role:    choice.Message.Role,
				Content: choice.Message.Content,
			},
			FinishReason: string(choice.FinishReason),
		}
	}

	response := models.ChatCompletionResponse{
		ID:      resp.ID,
		Object:  resp.Object,
		Created: resp.Created,
		Model:   resp.Model,
		Choices: choices,
		Usage: models.Usage{
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

