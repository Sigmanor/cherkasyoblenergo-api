package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadConfig_Success(t *testing.T) {
	viper.Reset()

	envVars := []string{"DB_NAME", "SERVER_PORT"}
	oldEnv := make(map[string]string)
	for _, key := range envVars {
		oldEnv[key] = os.Getenv(key)
		os.Unsetenv(key)
	}
	defer func() {
		for key, val := range oldEnv {
			if val != "" {
				os.Setenv(key, val)
			}
		}
		viper.Reset()
	}()

	tempDir, err := os.MkdirTemp("", "configtest")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	envContent := []byte(`DB_NAME=test.db
SERVER_PORT=8080
RATE_LIMIT_PER_MINUTE=30
CACHE_TTL_SECONDS=120
`)
	envFile := filepath.Join(tempDir, ".env")
	if err := os.WriteFile(envFile, envContent, 0o644); err != nil {
		t.Fatalf("failed to write .env file: %v", err)
	}

	cfg, err := LoadConfig(tempDir)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.DBName != "test.db" {
		t.Errorf("expected DBName 'test.db', got '%s'", cfg.DBName)
	}
	if cfg.ServerPort != "8080" {
		t.Errorf("expected SERVER_PORT '8080', got '%s'", cfg.ServerPort)
	}
	if cfg.RateLimitPerMinute != 30 {
		t.Errorf("expected RateLimitPerMinute 30, got %d", cfg.RateLimitPerMinute)
	}
	if cfg.CacheTTLSeconds != 120 {
		t.Errorf("expected CacheTTLSeconds 120, got %d", cfg.CacheTTLSeconds)
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	viper.Reset()

	envVars := []string{"DB_NAME", "SERVER_PORT", "RATE_LIMIT_PER_MINUTE", "CACHE_TTL_SECONDS"}
	oldEnv := make(map[string]string)
	for _, key := range envVars {
		oldEnv[key] = os.Getenv(key)
		os.Unsetenv(key)
	}
	defer func() {
		for key, val := range oldEnv {
			if val != "" {
				os.Setenv(key, val)
			}
		}
		viper.Reset()
	}()

	tempDir, err := os.MkdirTemp("", "configtest")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	envContent := []byte(`DB_NAME=test.db
SERVER_PORT=8080
`)
	envFile := filepath.Join(tempDir, ".env")
	if err := os.WriteFile(envFile, envContent, 0o644); err != nil {
		t.Fatalf("failed to write .env file: %v", err)
	}

	cfg, err := LoadConfig(tempDir)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.RateLimitPerMinute != 60 {
		t.Errorf("expected default RateLimitPerMinute 60, got %d", cfg.RateLimitPerMinute)
	}
	if cfg.CacheTTLSeconds != 60 {
		t.Errorf("expected default CacheTTLSeconds 60, got %d", cfg.CacheTTLSeconds)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

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
