package database

import (
	"fmt"

	"cherkasyoblenergo-api/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func createDSN(cfg config.Config) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
}

func ConnectDB(cfg config.Config) (*gorm.DB, error) {
	dsn := createDSN(cfg)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
