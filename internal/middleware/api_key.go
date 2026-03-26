package middleware

import (
	"crypto/subtle"

	"github.com/gofiber/fiber/v2"
)

func APIKeyAuth(apiKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if apiKey == "" {
			return c.Next()
		}

		providedKey := c.Get("X-API-Key")
		if providedKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "API key is required",
			})
		}

		if len(providedKey) != len(apiKey) {
			fakeKey := make([]byte, len(apiKey))
			subtle.ConstantTimeCompare([]byte(providedKey), fakeKey)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API key",
			})
		}

		if subtle.ConstantTimeCompare([]byte(providedKey), []byte(apiKey)) != 1 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API key",
			})
		}

		return c.Next()
	}
}
