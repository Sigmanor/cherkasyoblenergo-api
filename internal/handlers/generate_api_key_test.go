package handlers

import (
	"encoding/json"
	"net/http/httptest"
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

func TestGenerateAPIKey_Unauthorized(t *testing.T) {
	db := setupGenerateDB()
	cfg := config.Config{AdminPassword: "admin"}
	handler := GenerateAPIKey(db, cfg)
	app := fiber.New()
	app.Get("/generate", handler)

	req := httptest.NewRequest("GET", "/generate?admin_password=wrong", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
	}
}

func TestGenerateAPIKey_Success(t *testing.T) {
	db := setupGenerateDB()
	cfg := config.Config{AdminPassword: "admin"}
	handler := GenerateAPIKey(db, cfg)
	app := fiber.New()
	app.Get("/generate", handler)

	req := httptest.NewRequest("GET", "/generate?admin_password=admin&rate_limit=3", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	var result map[string]interface{}
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
