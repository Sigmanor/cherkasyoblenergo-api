package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"cherkasyoblenergo-api/internal/config"
	"cherkasyoblenergo-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupGenerateDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.APIKey{})
	return db
}

func TestCreateAPIKey_Unauthorized(t *testing.T) {
	db := setupGenerateDB()
	cfg := config.Config{AdminPassword: "admin"}
	handler := CreateAPIKey(db, cfg)
	app := fiber.New()
	app.Post("/api-keys", handler)

	req := httptest.NewRequest("POST", "/api-keys", strings.NewReader(`{"admin_password":"wrong"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
	}
}

func TestCreateAPIKey_SuccessWithJSON(t *testing.T) {
	db := setupGenerateDB()
	cfg := config.Config{AdminPassword: "admin"}
	handler := CreateAPIKey(db, cfg)
	app := fiber.New()
	app.Post("/api-keys", handler)

	req := httptest.NewRequest("POST", "/api-keys", strings.NewReader(`{"admin_password":"admin","rate_limit":3}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	if rateLimit, ok := result["rate_limit"]; ok {
		if int(rateLimit.(float64)) != 3 {
			t.Errorf("Expected rate_limit 3, got %v", rateLimit)
		}
	} else {
		t.Error("Expected rate_limit in response")
	}
	if _, ok := result["api_key"]; !ok {
		t.Error("Expected api_key in response")
	}
}

func TestCreateAPIKey_DefaultRateLimit(t *testing.T) {
	db := setupGenerateDB()
	cfg := config.Config{AdminPassword: "admin"}
	handler := CreateAPIKey(db, cfg)
	app := fiber.New()
	app.Post("/api-keys", handler)

	req := httptest.NewRequest("POST", "/api-keys", strings.NewReader(`{"admin_password":"admin"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	if rateLimit, ok := result["rate_limit"]; ok {
		if int(rateLimit.(float64)) != 6 {
			t.Errorf("Expected default rate_limit 6, got %v", rateLimit)
		}
	} else {
		t.Error("Expected rate_limit in response")
	}
}
