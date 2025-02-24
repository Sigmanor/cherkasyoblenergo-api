package handlers

import (
	"cherkasyoblenergo-api/internal/config"
	"cherkasyoblenergo-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func ManageAPIKey(db *gorm.DB, cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		adminPassword := c.Query("admin_password")
		if adminPassword != cfg.AdminPassword {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		key := c.Query("key")
		if key == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "API key is required",
			})
		}

		var apiKey models.APIKey
		if err := db.Where("key = ?", key).First(&apiKey).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "API key not found",
			})
		}

		if c.Query("update_key") == "true" {
			newKey := uuid.New().String()
			if err := db.Model(&apiKey).Update("key", newKey).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to update API key",
				})
			}
			return c.JSON(fiber.Map{
				"message": "API key updated successfully",
				"new_key": newKey,
			})
		}

		if c.Query("delete_key") == "true" {
			if err := db.Delete(&apiKey).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to delete API key",
				})
			}
			return c.JSON(fiber.Map{
				"message": "API key deleted successfully",
			})
		}

		if rateLimit := c.Query("update_rate_limit"); rateLimit != "" {
			if err := db.Model(&apiKey).Update("rate_limit", rateLimit).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to update rate limit",
				})
			}
			return c.JSON(fiber.Map{
				"message":        "Rate limit updated successfully",
				"new_rate_limit": rateLimit,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No valid operation specified",
		})
	}
}
