package middlewares

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"time"
)

func ApiKeyMiddleware(DB *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Getting the Authorization key a.k.a the APIKey
		apikey := c.Get("Authorization")
		slog.Info("APIkey: " + apikey)

		// Transaction stuffs
		tx, err := DB.Begin()
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()
		slog.Info("Begin transaction for APIKEY middleware")

		// Checks the APIKey's Integrity and Expiry
		var expiry time.Time

		// Query the expiry field along with other columns
		err = tx.QueryRowContext(c.Context(), "SELECT expires_at FROM apikeys WHERE token = ?", apikey).Scan(&expiry)
		if err != nil {
			if err == sql.ErrNoRows {
				slog.Error("API Key not found")
				return fiber.NewError(fiber.StatusNotFound, "API Key not found")
			}
			slog.Error("Error while querying API Key", "err", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
		}

		// Check if the token is expired
		now := time.Now()
		if now.After(expiry) {
			slog.Warn("API Key has expired", "expiry", expiry, "now", now)
			return fiber.NewError(fiber.StatusUnauthorized, "API Key has expired")
		}

		return c.Next()
	}
}
