package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wmbryce/agent-c/app/types"
)

func GetModels(c *fiber.Ctx) error {
	// models, err := store.GetModels()
	// if err != nil {
	// 	return nil, err
	// }
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"models": []types.Model{},
	})
}