package httpfiber

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/example/iching-fiber-app/internal/config"
	"github.com/example/iching-fiber-app/internal/domain"
	"github.com/example/iching-fiber-app/internal/service"
)

func NewApp(cfg config.Config, svc *service.ReadingService, embeddedStatic embed.FS) *fiber.App {
	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/health", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"status": "ok"}) })
	app.Get("/api/readings", func(c *fiber.Ctx) error { items, err := svc.List(c.Context()); if err != nil { return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()}) }; return c.JSON(items) })
	app.Get("/api/readings/:id", func(c *fiber.Ctx) error { item, err := svc.Get(c.Context(), c.Params("id")); if err != nil { return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()}) }; return c.JSON(item) })
	app.Post("/api/readings", func(c *fiber.Ctx) error { var in service.CreateReadingInput; if err := c.BodyParser(&in); err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()}) }; if in.Method == "" { in.Method = domain.MethodManual }; item, err := svc.Create(c.Context(), in); if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()}) }; return c.Status(fiber.StatusCreated).JSON(item) })
	app.Post("/api/readings/random", func(c *fiber.Ctx) error { var in service.CreateReadingInput; if err := c.BodyParser(&in); err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()}) }; if in.Question == "" { in.Question = "Random reading" }; in.Random = true; in.Method = domain.MethodCoins; item, err := svc.Create(c.Context(), in); if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()}) }; return c.Status(fiber.StatusCreated).JSON(item) })

	staticFS, err := fs.Sub(embeddedStatic, "static")
	if err != nil { panic(err) }
	app.Use("/", filesystem.New(filesystem.Config{Root: http.FS(staticFS), Index: "index.html", NotFoundFile: "index.html"}))
	return app
}
