package database

import (
	"fmt"
	"log"
	"strings"
	"time"

	"cherkasyoblenergo-api/internal/config"
	"cherkasyoblenergo-api/internal/logger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const defaultDBName = "cherkasyoblenergo.db"

func ConnectDB(cfg config.Config) (*gorm.DB, error) {
	dbPath := cfg.DBName
	if dbPath == "" {
		dbPath = defaultDBName
	}

	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=5000&_loc=Local", dbPath)
	dialector := sqlite.Open(dsn)
	log.Printf("Using SQLite database at: %s", dbPath)

	var gormLogLevel gormlogger.LogLevel
	switch strings.ToLower(cfg.LogLevel) {
	case "silent":
		gormLogLevel = gormlogger.Silent
	case "error":
		gormLogLevel = gormlogger.Error
	case "warn", "warning":
		gormLogLevel = gormlogger.Warn
	case "info":
		gormLogLevel = gormlogger.Warn
	case "debug":
		gormLogLevel = gormlogger.Info
	default:
		gormLogLevel = gormlogger.Info
	}

	newLogger := logger.NewGormLogger(gormLogLevel, time.Second)

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		log.Printf("Failed to enable foreign keys for SQLite: %v", err)
	}

	log.Printf("Successfully connected to database")
	return db, nil
}
