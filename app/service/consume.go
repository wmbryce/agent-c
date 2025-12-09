package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wmbryce/agent-c/app/types"
	"github.com/wmbryce/agent-c/app/utils"
)

// ConsumeModel func sends a request to the AI model provider.
// @Description Send a consume model request to the AI provider.
// @Summary consume an AI model
// @Tags AI
// @Accept json
// @Produce json
// @Param request body types.ConsumeModelRequest true "Consume model request"
// @Success 200 {object} types.ChatCompletionResponse
// @Security ApiKeyAuth
// @Router /v1/ai/consume [post]
func (s *Service) ConsumeModel(c *fiber.Ctx) error {
	request := &types.ConsumeModelRequest{}
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Get model credentials (endpoint URL and API key) in one query
	creds, err := s.store.GetModelCredentials(request.ModelKey)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "model not found or no API key available",
		})
	}

	// Check if tokens available cover the max cost
	if float64(creds.TokensAvailable) < request.MaxCost {
		return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
			"error": true,
			"msg":   "insufficient tokens available",
		})
	}

	// Build the request payload for the provider
	providerRequest := map[string]interface{}{
		"model":    request.ModelKey,
		"messages": request.Messages,
	}
	for k, v := range request.Options {
		providerRequest[k] = v
	}

	payload, err := json.Marshal(providerRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "failed to marshal request",
		})
	}

	// Send request to the model provider
	httpReq, err := http.NewRequest("POST", creds.RequestURL, bytes.NewBuffer(payload))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "failed to create request",
		})
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+creds.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": true,
			"msg":   "failed to reach model provider",
		})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "failed to read response",
		})
	}

	// Parse the provider response
	var response types.ChatCompletionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "failed to parse provider response",
		})
	}

	return c.JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"response": response,
	})
}
