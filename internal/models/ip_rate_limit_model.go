package models

import "time"

type IPRateLimit struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	IP           string    `gorm:"uniqueIndex;size:45;not null"`
	RequestCount int       `gorm:"not null;default:0"`
	WindowStart  time.Time `gorm:"not null"`
}
