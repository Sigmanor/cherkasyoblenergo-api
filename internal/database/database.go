package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"cherkasyoblenergo-api/internal/config"
	"cherkasyoblenergo-api/internal/logger"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const defaultDBName = "cherkasyoblenergo"

func getDBName(cfg config.Config) string {
	if cfg.DBName != "" {
		return cfg.DBName
	}
	return defaultDBName
}

func createDSN(cfg config.Config, dbName string) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, dbName, cfg.DBPort)
}

func ensureDatabaseExists(cfg config.Config) error {
	targetDBName := getDBName(cfg)
	dsn := createDSN(cfg, "template1")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres server: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping postgres server: %w", err)
	}

	var exists bool
	query := "SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)"
	err = db.QueryRow(query, targetDBName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		log.Printf("Database '%s' does not exist, creating it...", targetDBName)
		createQuery := fmt.Sprintf(`CREATE DATABASE "%s"`, targetDBName)
		_, err = db.Exec(createQuery)
		if err != nil {
			return fmt.Errorf("failed to create database '%s': %w", targetDBName, err)
		}
		log.Printf("Database '%s' created successfully", targetDBName)
	} else {
		log.Printf("Database '%s' already exists", targetDBName)
	}

	return nil
}

func ConnectDB(cfg config.Config) (*gorm.DB, error) {
	if err := ensureDatabaseExists(cfg); err != nil {
		return nil, fmt.Errorf("failed to ensure database exists: %w", err)
	}

	targetDBName := getDBName(cfg)
	dsn := createDSN(cfg, targetDBName)

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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database '%s': %w", targetDBName, err)
	}

	log.Printf("Successfully connected to database '%s'", targetDBName)
	return db, nil
}
