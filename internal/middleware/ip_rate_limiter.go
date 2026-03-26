package middleware

import (
	"log"
	"strings"
	"time"

	"cherkasyoblenergo-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type IPRateLimiter struct {
	db              *gorm.DB
	maxRequests     int
	windowDuration  time.Duration
	cleanupInterval time.Duration
	stopChan        chan struct{}
}

func NewIPRateLimiter(db *gorm.DB, maxRequests int) *IPRateLimiter {
	limiter := &IPRateLimiter{
		db:              db,
		maxRequests:     maxRequests,
		windowDuration:  time.Minute,
		cleanupInterval: 5 * time.Minute,
		stopChan:        make(chan struct{}),
	}

	go limiter.cleanupLoop()

	return limiter
}

func (rl *IPRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.cleanup()
		case <-rl.stopChan:
			return
		}
	}
}

func (rl *IPRateLimiter) cleanup() {
	cutoff := time.Now().Add(-rl.windowDuration * 2)
	result := rl.db.Where("window_start < ?", cutoff).Delete(&models.IPRateLimit{})
	if result.Error != nil {
		log.Printf("Error cleaning up rate limit entries: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("Cleaned up %d expired rate limit entries", result.RowsAffected)
	}
}

func (rl *IPRateLimiter) Stop() {
	close(rl.stopChan)
}

func (rl *IPRateLimiter) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		if ip == "" {
			ip = "unknown"
		}

		now := time.Now()
		windowStart := now.Truncate(rl.windowDuration)

		var currentCount int

		maxRetries := 3
		var err error
		for attempt := 0; attempt < maxRetries; attempt++ {
			err = rl.db.Transaction(func(tx *gorm.DB) error {
				var entry models.IPRateLimit
				result := tx.Where("ip = ?", ip).First(&entry)

				if result.Error == gorm.ErrRecordNotFound {
					entry = models.IPRateLimit{
						IP:           ip,
						RequestCount: 1,
						WindowStart:  windowStart,
					}
					if err := tx.Create(&entry).Error; err != nil {
						return err
					}
					currentCount = 1
					return nil
				}

				if result.Error != nil {
					return result.Error
				}

				if entry.WindowStart.Before(windowStart) {
					entry.WindowStart = windowStart
					entry.RequestCount = 1
				} else {
					entry.RequestCount++
				}

				currentCount = entry.RequestCount

				return tx.Save(&entry).Error
			})

			if err == nil || !strings.Contains(err.Error(), "database is locked") {
				break
			}
			time.Sleep(time.Millisecond * 50 * time.Duration(attempt+1))
		}

		if err != nil {
			log.Printf("Error in rate limit transaction: %v", err)
			return c.Next()
		}

		if currentCount > rl.maxRequests {
			c.Set("X-RateLimit-Limit", intToString(rl.maxRequests))
			c.Set("X-RateLimit-Remaining", "0")
			c.Set("Retry-After", "60")

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded. Please try again later.",
			})
		}

		remaining := rl.maxRequests - currentCount
		if remaining < 0 {
			remaining = 0
		}

		c.Set("X-RateLimit-Limit", intToString(rl.maxRequests))
		c.Set("X-RateLimit-Remaining", intToString(remaining))

		return c.Next()
	}
}

func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}
