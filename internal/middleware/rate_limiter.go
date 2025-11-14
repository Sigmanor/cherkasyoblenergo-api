package middleware

import (
	"context"
	"log"
	"sync"
	"time"

	"cherkasyoblenergo-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/ulule/limiter/v3"
	mem "github.com/ulule/limiter/v3/drivers/store/memory"
	"gorm.io/gorm"
)

var limiterCache = sync.Map{}

var skipPaths = map[string]struct{}{
	"/cherkasyoblenergo/api/generate-api-key": {},
	"/cherkasyoblenergo/api/update-api-key":   {},
	"/cherkasyoblenergo/api/api-keys":         {},
}

func RateLimiter(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if _, ok := skipPaths[c.Path()]; ok {
			return c.Next()
		}
		apiKey, ok := c.Locals("api_key").(models.APIKey)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get API key from context"})
		}
		lmt, ok := limiterCache.Load(apiKey.Key)
		if !ok {
			rate := limiter.Rate{Period: 1 * time.Minute, Limit: int64(apiKey.RateLimit)}
			store := mem.NewStore()
			lmt = limiter.New(store, rate)
			limiterCache.Store(apiKey.Key, lmt)
			log.Printf("Created new limiter for API key: %s with rate limit: %d", maskAPIKey(apiKey.Key), apiKey.RateLimit)
		} else {
			log.Printf("Using cached limiter for API key: %s", maskAPIKey(apiKey.Key))
		}
		limiterInstance := lmt.(*limiter.Limiter)
		var updatedAPIKey models.APIKey
		if err := db.Where("key = ?", apiKey.Key).First(&updatedAPIKey).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch API key from database"})
		}
		if updatedAPIKey.RateLimit != apiKey.RateLimit {
			rate := limiter.Rate{Period: 1 * time.Minute, Limit: int64(updatedAPIKey.RateLimit)}
			store := mem.NewStore()
			newLimiter := limiter.New(store, rate)
			limiterCache.Store(apiKey.Key, newLimiter)
			limiterInstance = newLimiter
			log.Printf("Updated limiter for API key: %s with new rate limit: %d", maskAPIKey(apiKey.Key), updatedAPIKey.RateLimit)
		}
		ctx := context.Background()
		limiterContext, err := limiterInstance.Get(ctx, apiKey.Key)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check rate limit"})
		}
		if limiterContext.Reached {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "You have reached maximum request limit."})
		}
		return c.Next()
	}
}
