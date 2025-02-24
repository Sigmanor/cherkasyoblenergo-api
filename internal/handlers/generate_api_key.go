package handlers

import (
	"strconv"

	"cherkasyoblenergo-api/internal/config"
	"cherkasyoblenergo-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GenerateAPIKey(db *gorm.DB, cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		adminPassword := c.Query("admin_password")
		if adminPassword != cfg.AdminPassword {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		rateLimitStr := c.Query("rate_limit", "2")
		rateLimit, err := strconv.Atoi(rateLimitStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid rate_limit value",
			})
		}

		key := uuid.New().String()

		apiKey := models.APIKey{
			Key:       key,
			RateLimit: rateLimit,
		}
		if err := db.Create(&apiKey).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create API key",
			})
		}

		return c.JSON(fiber.Map{
			"api_key":    key,
			"rate_limit": rateLimit,
		})
	}
}
