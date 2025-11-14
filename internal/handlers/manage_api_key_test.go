package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
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

func TestUpdateAPIKey_Unauthorized(t *testing.T) {
	db := setupManageDB()
	cfg := config.Config{AdminPassword: "admin"}
	handler := UpdateAPIKey(db, cfg)
	app := fiber.New()
	app.Patch("/api-keys", handler)

	body := bytes.NewBufferString(`{"admin_password":"wrong","key":"test-key","rotate_key": true}`)
	req := httptest.NewRequest("PATCH", "/api-keys", body)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
	}
}

func TestUpdateAPIKey_NotFound(t *testing.T) {
	db := setupManageDB()
	cfg := config.Config{AdminPassword: "admin"}
	handler := UpdateAPIKey(db, cfg)
	app := fiber.New()
	app.Patch("/api-keys", handler)

	body := bytes.NewBufferString(`{"admin_password":"admin","key":"missing","rotate_key": true}`)
	req := httptest.NewRequest("PATCH", "/api-keys", body)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("Expected status %d, got %d", fiber.StatusNotFound, resp.StatusCode)
	}
}

func TestUpdateAPIKey_RotateAndRateLimit(t *testing.T) {
	db := setupManageDB()
	cfg := config.Config{AdminPassword: "admin"}
	handler := UpdateAPIKey(db, cfg)
	app := fiber.New()
	app.Patch("/api-keys", handler)

	body := bytes.NewBufferString(`{"admin_password":"admin","key":"test-key","rotate_key": true, "rate_limit": 5}`)
	req := httptest.NewRequest("PATCH", "/api-keys", body)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	if _, ok := result["new_key"]; !ok {
		t.Error("Expected new_key in response")
	}
	if rateLimit, ok := result["new_rate_limit"]; ok {
		if int(rateLimit.(float64)) != 5 {
			t.Errorf("Expected new_rate_limit 5, got %v", rateLimit)
		}
	} else {
		t.Error("Expected new_rate_limit in response")
	}
}

func TestUpdateAPIKey_InvalidPayload(t *testing.T) {
	db := setupManageDB()
	cfg := config.Config{AdminPassword: "admin"}
	handler := UpdateAPIKey(db, cfg)
	app := fiber.New()
	app.Patch("/api-keys", handler)

	req := httptest.NewRequest("PATCH", "/api-keys", bytes.NewBufferString(`{"admin_password":"admin","key":"test-key"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}

func TestDeleteAPIKey_Success(t *testing.T) {
	db := setupManageDB()
	cfg := config.Config{AdminPassword: "admin"}
	handler := DeleteAPIKey(db, cfg)
	app := fiber.New()
	app.Delete("/api-keys", handler)

	req := httptest.NewRequest("DELETE", "/api-keys", bytes.NewBufferString(`{"admin_password":"admin","key":"test-key"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}

func TestDeleteAPIKey_NotFound(t *testing.T) {
	db := setupManageDB()
	cfg := config.Config{AdminPassword: "admin"}
	handler := DeleteAPIKey(db, cfg)
	app := fiber.New()
	app.Delete("/api-keys", handler)

	req := httptest.NewRequest("DELETE", "/api-keys", bytes.NewBufferString(`{"admin_password":"admin","key":"missing"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("Expected status %d, got %d", fiber.StatusNotFound, resp.StatusCode)
	}
}
