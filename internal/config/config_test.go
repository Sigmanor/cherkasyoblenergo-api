package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "configtest")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	envContent := []byte(`DB_HOST=localhost
DB_PORT=5432
DB_USER=testuser
DB_PASSWORD=testpass
DB_NAME=testdb
ADMIN_PASSWORD=adminpass
SERVER_PORT=8080
`)
	envFile := filepath.Join(tempDir, ".env")
	if err := os.WriteFile(envFile, envContent, 0644); err != nil {
		t.Fatalf("failed to write .env file: %v", err)
	}

	cfg, err := LoadConfig(tempDir)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.DBHost != "localhost" {
		t.Errorf("expected DBHost 'localhost', got '%s'", cfg.DBHost)
	}
	if cfg.DBPort != "5432" {
		t.Errorf("expected DBPort '5432', got '%s'", cfg.DBPort)
	}
	if cfg.DBUser != "testuser" {
		t.Errorf("expected DBUser 'testuser', got '%s'", cfg.DBUser)
	}
	if cfg.DBPassword != "testpass" {
		t.Errorf("expected DBPassword 'testpass', got '%s'", cfg.DBPassword)
	}
	if cfg.DBName != "testdb" {
		t.Errorf("expected DBName 'testdb', got '%s'", cfg.DBName)
	}
	if cfg.AdminPassword != "adminpass" {
		t.Errorf("expected AdminPassword 'adminpass', got '%s'", cfg.AdminPassword)
	}
	if cfg.SERVER_PORT != "8080" {
		t.Errorf("expected SERVER_PORT '8080', got '%s'", cfg.SERVER_PORT)
	}
	expectedVersion := "dev"
	if cfg.APP_VERSION != expectedVersion {
		t.Errorf("expected APP_VERSION '%s', got '%s'", expectedVersion, cfg.APP_VERSION)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "configtest")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	_, err = LoadConfig(tempDir)
	if err == nil {
		t.Errorf("expected error when .env file is missing, got nil")
	}
}
