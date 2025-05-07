package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"jirbthagoras/saas-starter/exceptions"
	"jirbthagoras/saas-starter/middlewares"
	"jirbthagoras/saas-starter/utils"
	"log"
	"log/slog"
	"time"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		slog.Error("Error loading .env file", "err", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: exceptions.ErrorHandler,
	})

	// Initialize database and Redis
	rdb := utils.NewRedisClient()
	db := utils.GetConnection()

	// Group Level Middleware - Apply Rate Limiter to all API routes
	api := app.Group("/api/v1", middlewares.RateLimiterMiddleware(rdb, 10, 60*time.Second))

	// PUBLIC ROUTES (No API Key Middleware)
	public := api.Group("/public")
	public.Get("/check", func(c *fiber.Ctx) error {
		return c.SendString("Public route - No API Key required")
	})
	public.Post("/api-key", CreateApiKeyHandler(db))

	// PROTECTED ROUTES (Requires API Key Middleware)
	protected := api.Group("/protected", middlewares.ApiKeyMiddleware(db))
	protected.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Protected route - API Key verified")
	})

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
