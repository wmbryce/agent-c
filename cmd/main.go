package main

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wmbryce/agent-c/app/middleware"
	"github.com/wmbryce/agent-c/app/routes"
	"github.com/wmbryce/agent-c/app/service"
	"github.com/wmbryce/agent-c/app/store"
	"github.com/wmbryce/agent-c/app/utils"
	"github.com/wmbryce/agent-c/cmd/configs"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/wmbryce/agent-c/docs"
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
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	ctx := context.Background()
	sqlStore := store.NewSqlStore(ctx)
	defer sqlStore.Close()

	app := fiber.New(configs.FiberConfig())
	middleware.FiberMiddleware(app)

	svc := service.New(&logger, sqlStore, app)
	routes.New(svc).Setup(app)

	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
