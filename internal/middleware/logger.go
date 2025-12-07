package middleware

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func maskAPIKey(apiKey string) string {
	if len(apiKey) < 10 {
		return apiKey
	}
	return apiKey[:5] + strings.Repeat("*", len(apiKey)-10) + apiKey[len(apiKey)-5:]
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

		attrs := []any{
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", c.Response().StatusCode()),
			slog.Duration("duration", duration),
			slog.String("ip", c.IP()),
		}

		if apiKey != "" {
			attrs = append(attrs, slog.String("api_key", maskAPIKey(apiKey)))
		}

		slog.Info("HTTP request", attrs...)

		return err
	}
}
