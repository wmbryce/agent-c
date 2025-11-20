package routes

import (
	"github.com/create-go-app/fiber-go-template/app/controllers"
	"github.com/create-go-app/fiber-go-template/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for POST method:
	route.Post("/book", middleware.JWTProtected(), controllers.CreateBook)              // create a new book
	route.Post("/user/sign/out", middleware.JWTProtected(), controllers.UserSignOut)    // de-authorization user
	route.Post("/token/renew", middleware.JWTProtected(), controllers.RenewTokens)      // renew Access & Refresh tokens
	route.Post("/openai/chat", middleware.JWTProtected(), controllers.ChatCompletion)   // send chat completion request to OpenAI
	
	// Blockchain routes:
	route.Post("/blockchain/balance", middleware.JWTProtected(), controllers.GetBalance)              // get ETH balance
	route.Post("/blockchain/receipt", middleware.JWTProtected(), controllers.GetTransactionReceipt)   // get transaction receipt
	route.Post("/blockchain/send-ether", middleware.JWTProtected(), controllers.SendEther)            // send ETH
	route.Get("/blockchain/info", middleware.JWTProtected(), controllers.GetBlockInfo)                // get blockchain info

	// Routes for PUT method:
	route.Put("/book", middleware.JWTProtected(), controllers.UpdateBook) // update one book by ID

	// Routes for DELETE method:
	route.Delete("/book", middleware.JWTProtected(), controllers.DeleteBook) // delete one book by ID
}
