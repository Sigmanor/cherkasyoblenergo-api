package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&Schedule{})
	db.Create(&Schedule{ID: 1, Title: "Test schedule", Date: time.Now()})
	return db
}

func TestPostSchedule_InvalidJSON(t *testing.T) {
	db := setupTestDB()
	handler := PostSchedule(db)
	app := fiber.New()
	app.Post("/schedules", handler)

	req := httptest.NewRequest("POST", "/schedules", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}

func TestPostSchedule_AllOption(t *testing.T) {
	db := setupTestDB()
	handler := PostSchedule(db)
	app := fiber.New()
	app.Post("/schedules", handler)

	body, _ := json.Marshal(map[string]any{
		"option": "all",
	})
	req := httptest.NewRequest("POST", "/schedules", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}
