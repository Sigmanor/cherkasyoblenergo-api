package handlers

import "github.com/gofiber/fiber/v2"

func isAdmin(c *fiber.Ctx, adminPassword string) bool {
	return c.Query("admin_password") == adminPassword
}
