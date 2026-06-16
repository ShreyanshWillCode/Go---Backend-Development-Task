package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/shreyxnsh/anyx-user-api/internal/handler"
	"github.com/shreyxnsh/anyx-user-api/internal/middleware"
)

func Register(app *fiber.App, userHandler *handler.UserHandler) {

	// Apply global security & utility middleware
	app.Use(helmet.New())
	app.Use(cors.New())
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	}))
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":	"ok",
			"service":	"anyx-user-api",
		})
	})

	users := app.Group("/users")
	{
		users.Post("/", userHandler.CreateUser)
		users.Get("/", userHandler.ListUsers)
		users.Get("/:id", userHandler.GetUser)
		users.Put("/:id", userHandler.UpdateUser)
		users.Delete("/:id", userHandler.DeleteUser)
	}
}
