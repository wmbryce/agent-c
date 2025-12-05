package main

import (
	"log"
	"os"

	"github.com/wmbryce/agent-c/app/middleware"
	"github.com/wmbryce/agent-c/app/routes"
	"github.com/wmbryce/agent-c/app/store"
	"github.com/wmbryce/agent-c/app/store/postgres"
	"github.com/wmbryce/agent-c/app/utils"
	"github.com/wmbryce/agent-c/cmd/configs"

	"github.com/gofiber/fiber/v2"

	_ "github.com/wmbryce/agent-c/docs" // load API Docs files (Swagger)

	_ "github.com/joho/godotenv/autoload" // load .env file automatically
)

// @title API
// @version 1.0
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email your@mail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Connect to PostgreSQL
	pool, err := postgres.PostgresConnection()
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer pool.Close()

	// Initialize the global store
	store.Init(pool)

	// Define Fiber config.
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Middlewares.
	middleware.FiberMiddleware(app) // Register Fiber's middleware for app.

	// Routes.
	routes.SetupRoutes(app) // Register a route for API Docs (Swagger).

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
