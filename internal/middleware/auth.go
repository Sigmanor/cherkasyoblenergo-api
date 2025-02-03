package middleware

import (
	"cherkasyoblenergo-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func APIKeyAuth(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() == "/cherkasyoblenergo/api/generate-api-key" ||
			c.Path() == "/cherkasyoblenergo/api/update-api-key" {
			return c.Next()
		}

		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "API key is missing",
			})
		}

		var key models.APIKey
		if err := db.Unscoped().Where("key = ? AND deleted_at IS NULL", apiKey).First(&key).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API key",
			})
		}

		c.Locals("api_key", key)
		return c.Next()
	}
}
