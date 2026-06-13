package routes

import (
	"github.com/anish/backend-development-task/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, userHandler *handler.UserHandler) {
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})
app.Get("/", func(c *fiber.Ctx) error {
	return c.SendString("users-api running")
})
	app.Post("/users", userHandler.CreateUser)
	app.Get("/users", userHandler.ListUsers)
	app.Get("/users/:id", userHandler.GetUserByID)
	app.Put("/users/:id", userHandler.UpdateUser)
	app.Delete("/users/:id", userHandler.DeleteUser)
}
