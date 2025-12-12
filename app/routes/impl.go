package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wmbryce/agent-c/app/service"
	"github.com/yokeTH/gofiber-scalar/scalar/v2"
)

type Routes struct {
	service *service.Service
}

func New(svc *service.Service) *Routes {
	return &Routes{service: svc}
}

func (r *Routes) Setup(app *fiber.App) {
	v1 := app.Group("api/v1")
	v1.Get("/ai/models", r.service.GetModels)
	v1.Post("/ai/models", r.service.CreateModel)
	v1.Post("/ai/consume", r.service.ConsumeModel)
	app.Get("/docs/*", scalar.New())

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "sorry, endpoint is not found",
		})
	})
}
