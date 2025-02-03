package models

import "gorm.io/gorm"

type APIKey struct {
	gorm.Model
	Key       string `gorm:"uniqueIndex" json:"key"`
	RateLimit int    `json:"rate_limit"`
}
