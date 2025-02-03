package middleware

import (
	"cherkasyoblenergo-api/internal/models"
	"context"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"gorm.io/gorm"
)

var limiterCache = sync.Map{}

func RateLimiter(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() == "/cherkasyoblenergo/api/generate-api-key" ||
			c.Path() == "/cherkasyoblenergo/api/update-api-key" {
			return c.Next()
		}

		apiKey, ok := c.Locals("api_key").(models.APIKey)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get API key from context",
			})
		}

		lmt, ok := limiterCache.Load(apiKey.Key)
		if !ok {
			rate := limiter.Rate{
				Period: 1 * time.Minute,
				Limit:  int64(apiKey.RateLimit),
			}
			store := memory.NewStore()
			lmt = limiter.New(store, rate)
			limiterCache.Store(apiKey.Key, lmt)
			log.Printf("Created new limiter for API key: %s with rate limit: %d", maskAPIKey(apiKey.Key), apiKey.RateLimit)
		} else {
			log.Printf("Using cached limiter for API key: %s", maskAPIKey(apiKey.Key))
		}

		limiterInstance := lmt.(*limiter.Limiter)

		var updatedAPIKey models.APIKey
		if err := db.Where("key = ?", apiKey.Key).First(&updatedAPIKey).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch API key from database",
			})
		}

		if updatedAPIKey.RateLimit != apiKey.RateLimit {
			rate := limiter.Rate{
				Period: 1 * time.Minute,
				Limit:  int64(updatedAPIKey.RateLimit),
			}
			store := memory.NewStore()
			newLimiter := limiter.New(store, rate)
			limiterCache.Store(apiKey.Key, newLimiter)
			log.Printf("Updated limiter for API key: %s with new rate limit: %d", maskAPIKey(apiKey.Key), updatedAPIKey.RateLimit)
		}

		ctx := context.Background()
		limiterContext, err := limiterInstance.Get(ctx, apiKey.Key)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check rate limit",
			})
		}

		if limiterContext.Reached {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "You have reached maximum request limit.",
			})
		}

		return c.Next()
	}
}
