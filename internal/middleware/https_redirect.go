package middleware

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func HTTPSEnforcement() fiber.Handler {
	forceHTTPS := os.Getenv("FORCE_HTTPS") == "true"

	return func(c *fiber.Ctx) error {
		if !forceHTTPS {
			return c.Next()
		}

		if c.Protocol() == "https" {
			return c.Next()
		}
		
		if c.Get("X-Forwarded-Proto") == "https" {
			return c.Next()
		}

		target := "https://" + c.Hostname() + c.OriginalURL()
		
		if strings.Contains(c.Hostname(), "localhost") {
			return c.Next() 
		}

		log.Printf("Redirecting to HTTPS: %s", target)
		return c.Redirect(target, fiber.StatusMovedPermanently)
	}
}
