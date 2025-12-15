package models

import "gorm.io/gorm"

type APIKey struct {
	gorm.Model
	Key                string  `gorm:"uniqueIndex" json:"key"`
	RateLimit          int     `json:"rate_limit"`
	WebhookURL         string  `gorm:"type:varchar(255)" json:"webhook_url"`
	WebhookEnabled     bool    `gorm:"default:false" json:"webhook_enabled"`
	WebhookFailedAttempts int  `gorm:"default:0" json:"webhook_failed_attempts"`
}
