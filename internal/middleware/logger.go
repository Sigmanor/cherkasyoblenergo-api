package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func maskAPIKey(apiKey string) string {
	if len(apiKey) < 10 {
		return apiKey
	}
	return fmt.Sprintf("%s%s%s", apiKey[:5], strings.Repeat("*", len(apiKey)-10), apiKey[len(apiKey)-5:])
}

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}

		maskedKey := ""
		if apiKey != "" {
			maskedKey = fmt.Sprintf(" [API Key: %s]", maskAPIKey(apiKey))
		}

		fmt.Printf("[%s] %s %s %d %v %s%s\n",
			time.Now().Format("2006-01-02 15:04:05"),
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
			duration,
			c.IP(),
			maskedKey,
		)

		return err
	}
}
