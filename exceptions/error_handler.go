package exceptions

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var fiberErr *fiber.Error
	if !errors.As(err, &fiberErr) {
		slog.Error("Internal Server Error")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"errors":  nil,
		})
	}

	slog.Error(fiberErr.Message)
	return c.Status(fiberErr.Code).JSON(fiber.Map{
		"message": fiberErr.Message,
		"errors":  nil,
	})
}
