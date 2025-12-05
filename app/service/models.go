package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wmbryce/agent-c/app/store"
	"github.com/wmbryce/agent-c/app/types"
	"github.com/wmbryce/agent-c/app/utils"
)

func GetModels(c *fiber.Ctx) error {
	models, err := store.GetModels()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":  false,
		"msg":    nil,
		"models": models,
	})
}

func CreateModel(c *fiber.Ctx) error {
	model := &types.Model{}
	if err := c.BodyParser(model); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(model); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	createdModel, err := store.CreateModel(model)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "Model created successfully",
		"model": createdModel,
	})
}
