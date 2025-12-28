package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// APIKeyAuth creates a middleware that optionally requires an API key.
// If apiKey is empty, all requests are allowed (public mode).
// If apiKey is set, requests must include a matching X-API-Key header.
func APIKeyAuth(apiKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// If no API key configured, allow all requests (public mode)
		if apiKey == "" {
			return c.Next()
		}

		// Check for API key in header
		providedKey := c.Get("X-API-Key")
		if providedKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "API key is required",
			})
		}

		// Validate the API key
		if providedKey != apiKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API key",
			})
		}

		return c.Next()
	}
}
