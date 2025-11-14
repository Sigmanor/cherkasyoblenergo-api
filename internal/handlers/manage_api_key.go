package handlers

import (
	"cherkasyoblenergo-api/internal/config"
	"cherkasyoblenergo-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type updateAPIKeyRequest struct {
	AdminPassword string `json:"admin_password"`
	Key           string `json:"key"`
	RotateKey     bool   `json:"rotate_key"`
	RateLimit     *int   `json:"rate_limit"`
}

func UpdateAPIKey(db *gorm.DB, cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req updateAPIKeyRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format"})
		}

		if !isAdminPasswordValid(req.AdminPassword, cfg.AdminPassword) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		if req.Key == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "API key is required"})
		}

		if !req.RotateKey && req.RateLimit == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "At least one update field must be specified"})
		}

		var apiKey models.APIKey
		if err := db.Where("key = ?", req.Key).First(&apiKey).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "API key not found"})
		}

		updates := make(map[string]interface{})
		response := fiber.Map{"message": "API key updated successfully"}

		if req.RotateKey {
			newKey := uuid.New().String()
			updates["key"] = newKey
			response["new_key"] = newKey
		}

		if req.RateLimit != nil {
			if *req.RateLimit <= 0 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "rate_limit must be greater than zero"})
			}
			updates["rate_limit"] = *req.RateLimit
			response["new_rate_limit"] = *req.RateLimit
		}

		if err := db.Model(&apiKey).Updates(updates).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update API key"})
		}

		return c.JSON(response)
	}
}

type deleteAPIKeyRequest struct {
	AdminPassword string `json:"admin_password"`
	Key           string `json:"key"`
}

func DeleteAPIKey(db *gorm.DB, cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req deleteAPIKeyRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format"})
		}

		if !isAdminPasswordValid(req.AdminPassword, cfg.AdminPassword) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		if req.Key == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "API key is required"})
		}

		var apiKey models.APIKey
		if err := db.Where("key = ?", req.Key).First(&apiKey).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "API key not found"})
		}

		if err := db.Delete(&apiKey).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete API key"})
		}

		return c.JSON(fiber.Map{"message": "API key deleted successfully"})
	}
}
