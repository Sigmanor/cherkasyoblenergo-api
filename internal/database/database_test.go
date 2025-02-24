package database

import (
	"cherkasyoblenergo-api/internal/config"
	"testing"
)

func TestConnectDB_InvalidConfig(t *testing.T) {
	cfg := config.Config{
		DBHost:     "invalid_host",
		DBUser:     "user",
		DBPassword: "pass",
		DBName:     "dbname",
		DBPort:     "5432",
	}
	_, err := ConnectDB(cfg)
	if err == nil {
		t.Error("Expected error with invalid config, got nil")
	}
}
