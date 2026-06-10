package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"maz/domain/shared"
)

// ErrorHandler maps unhandled Fiber errors to JSON responses.
func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "internal server error"

	var fe *fiber.Error
	if errors.As(err, &fe) {
		code = fe.Code
		msg = fe.Message
	}

	if errors.Is(err, shared.ErrNotFound) {
		code = fiber.StatusNotFound
		msg = err.Error()
	}

	return c.Status(code).JSON(errorBody{Error: msg})
}
