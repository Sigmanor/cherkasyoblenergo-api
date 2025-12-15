package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"

	"cherkasyoblenergo-api/internal/models"
)

func ValidateWebhookURL(url string) error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to reach webhook URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook URL returned status code %d, expected 200", resp.StatusCode)
	}

	return nil
}

func SendWebhook(apiKey models.APIKey, schedules []models.Schedule) error {
	if apiKey.WebhookURL == "" {
		return fmt.Errorf("webhook URL is empty")
	}

	payload := map[string]interface{}{
		"schedules": schedules,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("POST", apiKey.WebhookURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Token", apiKey.Key)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook returned status code %d, expected 200", resp.StatusCode)
	}

	return nil
}

func SendWebhookWithRetry(db *gorm.DB, apiKey models.APIKey, schedules []models.Schedule) {
	delays := []time.Duration{0 * time.Second, 60 * time.Second, 600 * time.Second}

	for attempt, delay := range delays {
		if delay > 0 {
			time.Sleep(delay)
		}

		var freshAPIKey models.APIKey
		if err := db.Where("id = ?", apiKey.ID).First(&freshAPIKey).Error; err != nil {
			fmt.Printf("Failed to reload API key %s from database: %v\n", apiKey.Key, err)
			return
		}
		apiKey = freshAPIKey

		err := SendWebhook(apiKey, schedules)
		if err == nil {
			if err := db.Model(&apiKey).Update("WebhookFailedAttempts", 0).Error; err != nil {
				fmt.Printf("Failed to reset webhook failed attempts for API key %s: %v\n", apiKey.Key, err)
			} else {
				fmt.Printf("Webhook sent successfully for API key %s after %d attempts\n", apiKey.Key, attempt+1)
			}
			return
		}

		newFailedAttempts := apiKey.WebhookFailedAttempts + 1
		if err := db.Model(&apiKey).Update("WebhookFailedAttempts", newFailedAttempts).Error; err != nil {
			fmt.Printf("Failed to update webhook failed attempts for API key %s: %v\n", apiKey.Key, err)
		} else {
			apiKey.WebhookFailedAttempts = newFailedAttempts
		}

		if newFailedAttempts >= 3 {
			if err := db.Model(&apiKey).Update("WebhookEnabled", false).Error; err != nil {
				fmt.Printf("Failed to disable webhook for API key %s: %v\n", apiKey.Key, err)
			} else {
				fmt.Printf("Webhook disabled for API key %s after %d failed attempts\n", apiKey.Key, newFailedAttempts)
			}
			return
		}

		fmt.Printf("Webhook attempt %d failed for API key %s: %v\n", attempt+1, apiKey.Key, err)
	}

	fmt.Printf("All webhook attempts failed for API key %s\n", apiKey.Key)
}

func TriggerWebhooks(db *gorm.DB, schedules []models.Schedule) {
	var apiKeys []models.APIKey

	result := db.Where("webhook_url IS NOT NULL AND webhook_url != '' AND webhook_enabled = true AND deleted_at IS NULL").Find(&apiKeys)
	if result.Error != nil {
		fmt.Printf("Failed to fetch API keys for webhooks: %v\n", result.Error)
		return
	}

	if len(apiKeys) == 0 {
		fmt.Println("No enabled API keys with webhook URLs found")
		return
	}

	fmt.Printf("Found %d enabled API keys with webhook URLs\n", len(apiKeys))

	for _, apiKey := range apiKeys {
		go func(key models.APIKey) {
			SendWebhookWithRetry(db, key, schedules)
		}(apiKey)
	}
}
