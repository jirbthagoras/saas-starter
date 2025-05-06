package middlewares

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
		key := c.IP()

		// Calling the rate limiter
		allow, count, err := utils.AllowRequest(rdb, key, maxRequest, window)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
		}

		// Checks the result
		if !allow {
			slog.Info("Request Blocked")
			return fiber.NewError(fiber.StatusUnauthorized, "Try again")
		}

		slog.Info(fmt.Sprintf("The key: %s, have %d requests left", key, count))
		return c.Next()
	}
}
