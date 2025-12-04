package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wmbryce/agent-c/app/types"
	"github.com/wmbryce/agent-c/app/utils"
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
// @Router /v1/ai/chat [post]
func AiChatCompletion(c *fiber.Ctx) error {
	// Create new ChatCompletionRequest struct
	request := &types.ChatCompletionRequest{}
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	response := types.ChatCompletionResponse{
		ID:      "123",
		Object:  "chat.completion",
		Created: 1717334400,
		Model:   "gpt-4o",
		Choices: []types.Choice{
			{Index: 0, Message: types.ChatMessage{Role: "assistant", Content: "Hello, how can I help you today?"}},
		},
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"response": response,
	})
}


