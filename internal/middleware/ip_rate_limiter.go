package middleware

import (
	"log"
	"sync"
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
	mu              sync.RWMutex
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

		rl.mu.Lock()

		var entry models.IPRateLimit
		err := rl.db.Where("ip = ?", ip).First(&entry).Error

		if err == gorm.ErrRecordNotFound {
			entry = models.IPRateLimit{
				IP:           ip,
				RequestCount: 1,
				WindowStart:  windowStart,
			}
			if err := rl.db.Create(&entry).Error; err != nil {
				rl.mu.Unlock()
				log.Printf("Error creating rate limit entry: %v", err)
				return c.Next()
			}
			rl.mu.Unlock()

			c.Set("X-RateLimit-Limit", intToString(rl.maxRequests))
			c.Set("X-RateLimit-Remaining", intToString(rl.maxRequests-1))
			return c.Next()
		}

		if err != nil {
			rl.mu.Unlock()
			log.Printf("Error fetching rate limit entry: %v", err)
			return c.Next()
		}

		if entry.WindowStart.Before(windowStart) {
			entry.RequestCount = 1
			entry.WindowStart = windowStart
		} else {
			entry.RequestCount++
		}

		if entry.RequestCount > rl.maxRequests {
			rl.db.Save(&entry)
			rl.mu.Unlock()

			c.Set("X-RateLimit-Limit", intToString(rl.maxRequests))
			c.Set("X-RateLimit-Remaining", "0")
			c.Set("Retry-After", "60")

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded. Please try again later.",
			})
		}

		rl.db.Save(&entry)
		rl.mu.Unlock()

		remaining := rl.maxRequests - entry.RequestCount
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
