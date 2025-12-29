package database

import (
	"os"
	"testing"

	"cherkasyoblenergo-api/internal/config"
)

func TestConnectDB_SQLite(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test_*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	cfg := config.Config{
		DBName:   tempFile.Name(),
		LogLevel: "silent",
	}

	db, err := ConnectDB(cfg)
	if err != nil {
		t.Fatalf("ConnectDB failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get underlying sql.DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Errorf("failed to ping database: %v", err)
	}

	sqlDB.Close()
}
