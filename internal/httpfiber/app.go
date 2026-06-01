package httpfiber

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"github.com/example/iching-fiber-app/internal/config"
	"github.com/example/iching-fiber-app/internal/service"
)

func NewApp(_ config.Config, svc *service.ReadingService) *fiber.App {
	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	app.Get("/api/readings", func(c fiber.Ctx) error {
		items, err := svc.List(c.Context())
		if err != nil { return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()}) }
		return c.JSON(items)
	})
	app.Get("/api/readings/:id", func(c fiber.Ctx) error {
		item, err := svc.Get(c.Context(), c.Params("id"))
		if err != nil { return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()}) }
		return c.JSON(item)
	})
	app.Post("/api/readings", func(c fiber.Ctx) error {
		var in service.CreateReadingInput
		if err := c.Bind().Body(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		item, err := svc.Create(c.Context(), in)
		if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()}) }
		return c.Status(fiber.StatusCreated).JSON(item)
	})

	app.Static("/", "./web/static")
	return app
}
