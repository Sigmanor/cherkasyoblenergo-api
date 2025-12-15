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

type limiterEntry struct {
	limiter     *limiter.Limiter
	rateLimit   int
	lastChecked time.Time
}

const rateLimitCacheTTL = 60 * time.Second

var limiterCache = sync.Map{}

var skipPaths = map[string]struct{}{
	"/cherkasyoblenergo/api/api-keys": {},
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

		entry, exists := limiterCache.Load(apiKey.Key)
		var limiterInstance *limiter.Limiter

		if !exists {
			rate := limiter.Rate{Period: 1 * time.Minute, Limit: int64(apiKey.RateLimit)}
			store := mem.NewStore()
			limiterInstance = limiter.New(store, rate)
			limiterCache.Store(apiKey.Key, &limiterEntry{
				limiter:     limiterInstance,
				rateLimit:   apiKey.RateLimit,
				lastChecked: time.Now(),
			})
			log.Printf("Created new limiter for API key: %s with rate limit: %d", maskAPIKey(apiKey.Key), apiKey.RateLimit)
		} else {
			cachedEntry := entry.(*limiterEntry)
			limiterInstance = cachedEntry.limiter

			if time.Since(cachedEntry.lastChecked) > rateLimitCacheTTL {
				var updatedAPIKey models.APIKey
				if err := db.Where("key = ?", apiKey.Key).First(&updatedAPIKey).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch API key from database"})
				}

				if updatedAPIKey.RateLimit != cachedEntry.rateLimit {
					rate := limiter.Rate{Period: 1 * time.Minute, Limit: int64(updatedAPIKey.RateLimit)}
					store := mem.NewStore()
					limiterInstance = limiter.New(store, rate)
					log.Printf("Updated limiter for API key: %s with new rate limit: %d", maskAPIKey(apiKey.Key), updatedAPIKey.RateLimit)
				}

				limiterCache.Store(apiKey.Key, &limiterEntry{
					limiter:     limiterInstance,
					rateLimit:   updatedAPIKey.RateLimit,
					lastChecked: time.Now(),
				})
			}
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
