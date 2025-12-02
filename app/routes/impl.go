package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wmbryce/agent-c/app/controllers"

	swagger "github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App) {
	v1 := app.Group("api/v1")
	v1.Post("/ai/chat", controllers.AiChatCompletion)

	// Routes for GET method:
	swaggerRoute := app.Group("/swagger")
	// Routes for GET method:
	swaggerRoute.Get("*", swagger.HandlerDefault) // get one user by ID


	app.Use(
		// Anonymous function.
		func(c *fiber.Ctx) error {
			// Return HTTP 404 status and JSON response.
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   "sorry, endpoint is not found",
			})
		},
	)
}