package main

import (
	"github.com/gofiber/fiber/v2"
	"jirbthagoras/saas-starter/exceptions"
)

func main() {
	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: exceptions.ErrorHandler,
	})
	api := app.Group("/api/v1")

	// Create basic path
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello World",
		})
	})
}
