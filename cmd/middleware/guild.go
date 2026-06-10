package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const guildHeader = "X-Guild-ID"
const guildIDLocalKey = "guild_id"

// RequireGuildID validates the X-Guild-ID header on write operations.
func RequireGuildID(c *fiber.Ctx) error {
	raw := c.Get(guildHeader)
	if raw == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "X-Guild-ID header is required",
		})
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid X-Guild-ID",
		})
	}
	c.Locals(guildIDLocalKey, id)
	return c.Next()
}

// GuildID returns the guild ID set by RequireGuildID.
func GuildID(c *fiber.Ctx) (uuid.UUID, bool) {
	id, ok := c.Locals(guildIDLocalKey).(uuid.UUID)
	return id, ok
}
