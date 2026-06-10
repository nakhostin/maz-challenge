package fiberx

import (
	"github.com/gofiber/fiber/v2"

	"maz/domain/shared"
)

// BodyParser decodes the request body into T using Fiber's parser.
func BodyParser[T any](c *fiber.Ctx) (T, error) {
	var body T
	if err := c.BodyParser(&body); err != nil {
		return body, shared.ErrInvalidState
	}
	return body, nil
}
