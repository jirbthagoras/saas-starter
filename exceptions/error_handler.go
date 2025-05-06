package exceptions

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var fiberErr *fiber.Error
	if !errors.As(err, &fiberErr) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"errors":  nil,
		})
	}

	return c.Status(fiberErr.Code).JSON(fiber.Map{
		"message": fiberErr.Message,
		"errors":  nil,
	})
}
