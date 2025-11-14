package handlers

import (
	"cherkasyoblenergo-api/internal/config"
	"cherkasyoblenergo-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type createAPIKeyRequest struct {
	AdminPassword string `json:"admin_password"`
	RateLimit     *int   `json:"rate_limit"`
}

func CreateAPIKey(db *gorm.DB, cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req createAPIKeyRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format"})
		}

		if !isAdminPasswordValid(req.AdminPassword, cfg.AdminPassword) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		rateLimit := 2
		if req.RateLimit != nil {
			if *req.RateLimit <= 0 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "rate_limit must be greater than zero"})
			}
			rateLimit = *req.RateLimit
		}

		key := uuid.New().String()
		apiKey := models.APIKey{Key: key, RateLimit: rateLimit}
		if err := db.Create(&apiKey).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create API key"})
		}
		return c.JSON(fiber.Map{"api_key": key, "rate_limit": rateLimit})
	}
}
