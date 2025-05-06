package main

import (
	"github.com/gofiber/fiber/v2"
	"jirbthagoras/saas-starter/exceptions"
	"jirbthagoras/saas-starter/middlewares"
	"jirbthagoras/saas-starter/utils"
	"time"
)

func main() {
	//Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: exceptions.ErrorHandler,
	})
	api := app.Group("/api/v1")
	rdb := utils.NewRedisClient()

	//Apply some Middlewares
	api.Use(middlewares.RateLimiterMiddleware(rdb, 10, 60*time.Second))

	//Create basic path
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello World",
		})
	})

	app.Listen(":3000")
}
