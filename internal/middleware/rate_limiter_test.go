package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cherkasyoblenergo-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type APIKey struct {
	ID        uint
	Key       string
	RateLimit int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func TestRateLimiter(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&APIKey{})
	testKey := APIKey{Key: "rate-limit-key", RateLimit: 2}
	db.Create(&testKey)

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("api_key", models.APIKey{
			Key:       testKey.Key,
			RateLimit: testKey.RateLimit,
		})
		return c.Next()
	})
	app.Use(RateLimiter(db))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 for request %d, got %d", i+1, resp.StatusCode)
		}
	}
	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusTooManyRequests {
		t.Errorf("Expected status 429 on third request, got %d", resp.StatusCode)
	}
}
