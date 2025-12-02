package routes

import (
	"github.com/create-go-app/fiber-go-template/app/controllers"
	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// route.Post("/ai/chat", controllers.ChatCompletion)   // send chat completion request to OpenAI
	route.Post("/ai/chat", controllers.AiChatCompletion)   // send chat completion request to OpenAI

}
