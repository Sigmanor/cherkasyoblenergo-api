package handlers

import (
	"cherkasyoblenergo-api/internal/middleware"
	"cherkasyoblenergo-api/internal/webhook"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type registerWebhookRequest struct {
	WebhookURL string `json:"webhook_url"`
}

type webhookStatusResponse struct {
	WebhookURL            string `json:"webhook_url"`
	WebhookEnabled        bool   `json:"webhook_enabled"`
	WebhookFailedAttempts int    `json:"webhook_failed_attempts"`
}

func RegisterWebhook(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey, err := middleware.GetAPIKeyFromContext(c)
		if err != nil {
			return err
		}

		var req registerWebhookRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format"})
		}

		if req.WebhookURL == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "webhook_url is required"})
		}

		if err := webhook.ValidateWebhookURL(req.WebhookURL); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid webhook URL: " + err.Error()})
		}

		updates := map[string]interface{}{
			"WebhookURL":            req.WebhookURL,
			"WebhookEnabled":        true,
			"WebhookFailedAttempts": 0,
		}

		if err := db.Model(apiKey).Updates(updates).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update webhook settings"})
		}

		return c.JSON(fiber.Map{"message": "Webhook registered successfully", "webhook_url": req.WebhookURL})
	}
}

func DeleteWebhook(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey, err := middleware.GetAPIKeyFromContext(c)
		if err != nil {
			return err
		}

		updates := map[string]interface{}{
			"WebhookURL":     "",
			"WebhookEnabled": false,
		}

		if err := db.Model(apiKey).Updates(updates).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete webhook"})
		}

		return c.JSON(fiber.Map{"message": "Webhook deleted successfully"})
	}
}

func GetWebhookStatus(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey, err := middleware.GetAPIKeyFromContext(c)
		if err != nil {
			return err
		}

		if err := db.First(apiKey, apiKey.ID).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch webhook status"})
		}

		response := webhookStatusResponse{
			WebhookURL:            apiKey.WebhookURL,
			WebhookEnabled:        apiKey.WebhookEnabled,
			WebhookFailedAttempts: apiKey.WebhookFailedAttempts,
		}

		return c.JSON(response)
	}
}
