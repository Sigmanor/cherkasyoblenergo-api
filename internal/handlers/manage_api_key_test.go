package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"strconv"
	"testing"

	"cherkasyoblenergo-api/internal/config"
	"cherkasyoblenergo-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupManageDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.APIKey{})
	db.Create(&models.APIKey{Key: "test-key", RateLimit: 2})
	return db
}

func TestManageAPIKey_Unauthorized(t *testing.T) {
	db := setupManageDB()
	cfg := config.Config{AdminPassword: "admin"}
	handler := ManageAPIKey(db, cfg)
	app := fiber.New()
	app.Get("/manage", handler)

	req := httptest.NewRequest("GET", "/manage?admin_password=wrong&key=test-key", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
	}
}

func TestManageAPIKey_NotFound(t *testing.T) {
	db := setupManageDB()
	cfg := config.Config{AdminPassword: "admin"}
	handler := ManageAPIKey(db, cfg)
	app := fiber.New()
	app.Get("/manage", handler)

	req := httptest.NewRequest("GET", "/manage?admin_password=admin&key=nonexistent", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("Expected status %d, got %d", fiber.StatusNotFound, resp.StatusCode)
	}
}

func TestManageAPIKey_UpdateRateLimit(t *testing.T) {
	db := setupManageDB()
	cfg := config.Config{AdminPassword: "admin"}
	handler := ManageAPIKey(db, cfg)
	app := fiber.New()
	app.Get("/manage", handler)

	req := httptest.NewRequest("GET", "/manage?admin_password=admin&key=test-key&update_rate_limit=5", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	if newRate, ok := result["new_rate_limit"]; ok {
		if strconv.Itoa(5) != newRate.(string) {
			t.Errorf("Expected new_rate_limit to be '5', got %v", newRate)
		}
	} else {
		t.Error("Expected new_rate_limit in response")
	}
}
