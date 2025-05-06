package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"jirbthagoras/saas-starter/utils"
	"log/slog"
	"time"
)

var ctx = context.Background()

func RateLimiterMiddleware(rdb *redis.Client, maxRequest int, window time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Getting the stuffs
		key := c.Get("Authorization")

		// Calling the rate limiter
		allow, count, err := utils.AllowRequest(rdb, key, 5, 60*time.Second)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
		}

		// Checks the result
		if !allow {
			slog.Info("Request Blocked yahahah")
			return fiber.NewError(fiber.StatusUnauthorized, "Try again")
		}

		slog.Info(fmt.Sprintf("The key: %s, have %d requests left", key, count))
		return c.Next()
	}
}
