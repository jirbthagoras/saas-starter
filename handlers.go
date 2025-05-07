package main

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"jirbthagoras/saas-starter/utils"
	"log/slog"
	"time"
)

func CreateApiKeyHandler(DB *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Transaction
		tx, err := DB.Begin()
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()

		if err != nil {
			slog.Error("Failed to begin transaction")
			return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
		}

		// Generate a random strings
		key, err := utils.GenerateRandomString(10)

		if err != nil {
			slog.Error("Failed to create token")
			return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
		}

		expiresAt := time.Now().AddDate(0, 0, 30)

		// Creates a new APIKey

		result, err := tx.ExecContext(c.Context(), "INSERT INTO `apikeys` (token, expires_at) VALUES(?, ?)", key, expiresAt)
		// The error handling is very important guys.
		if err != nil {
			slog.Error("Failed to execute query", "err", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to insert key")
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			slog.Info("No rows affected while inserting key", "err", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to insert key")
		}

		return c.JSON(fiber.Map{
			"message": "Success creating key",
			"key":     key,
		})
	}
}
